package gapi

import (
	"fmt"

	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/pb"
	"github.com/umarhadi/bank-server/token"
	"github.com/umarhadi/bank-server/util"
	"github.com/umarhadi/bank-server/worker"
)

type Server struct {
	pb.UnimplementedBankServerServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

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
