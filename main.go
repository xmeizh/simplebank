package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq" // Without this package our server cannot connect to the db
	"github.com/rakyll/statik/fs"
	"github.com/xmeizh/simplebank/api"
	db "github.com/xmeizh/simplebank/db/postgresql"
	_ "github.com/xmeizh/simplebank/doc/statik"
	"github.com/xmeizh/simplebank/gapi"
	"github.com/xmeizh/simplebank/pb"
	"github.com/xmeizh/simplebank/util"
	"github.com/xmeizh/simplebank/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		// enable human-friendly logging
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	// start redis worker
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		gapi.LogFatal("cannot create new migrate instance:", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		gapi.LogFatal("failed to run migrate up:", err)
	}
	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		gapi.LogFatal("failed to start task processor", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		gapi.LogFatal("cannot create server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		gapi.LogFatal("cannot start server", err)
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		gapi.LogFatal("cannot create server:", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		gapi.LogFatal("failed to listen", err)
	}
	log.Info().Msgf("start gRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		gapi.LogFatal("cannot start gRPC server", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		gapi.LogFatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	gwmux := runtime.NewServeMux(jsonOption)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, gwmux, server)
	if err != nil {
		gapi.LogFatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	statikFS, err := fs.New()
	if err != nil {
		gapi.LogFatal("cannot create statik fs", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		gapi.LogFatal("failed to listen", err)
	}

	log.Info().Msgf("start HTTP gateway server at %s", lis.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(lis, handler)
	if err != nil {
		gapi.LogFatal("cannot start HTTP gateway server", err)
	}
}
