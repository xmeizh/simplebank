package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Println(config.DBDriver)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runDBMigration(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("db migrated successfully")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("start gRPC server at %s", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
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
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik fs: %v", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("start HTTP gateway server at %s", lis.Addr().String())
	err = http.Serve(lis, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server")
	}
}
