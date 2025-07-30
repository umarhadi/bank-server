// Package gapi provides gRPC API implementation for the bank server.
// It includes user management, authentication, and other banking operations.
package gapi

import (
	"fmt"

	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/pb"
	"github.com/umarhadi/bank-server/token"
	"github.com/umarhadi/bank-server/util"
	"github.com/umarhadi/bank-server/worker"
)

// Server serves gRPC requests for banking services.
// It implements the BankServerServer interface generated from protobuf definitions.
type Server struct {
	pb.UnimplementedBankServerServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server instance with the provided dependencies.
// It initializes the token maker using PASETO tokens for enhanced security.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
