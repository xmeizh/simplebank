package gapi

import (
	"fmt"

	db "github.com/xmeizh/simplebank/db/postgresql"
	"github.com/xmeizh/simplebank/pb"
	"github.com/xmeizh/simplebank/token"
	"github.com/xmeizh/simplebank/util"
)

// Server serves HTTP requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
