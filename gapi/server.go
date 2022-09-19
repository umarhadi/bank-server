package gapi

import (
	"fmt"

	db "github.com/umarhadi/bank-server/db/sqlc"
	"github.com/umarhadi/bank-server/pb"
	"github.com/umarhadi/bank-server/token"
	"github.com/umarhadi/bank-server/util"
)

type Server struct {
	pb.UnimplementedBankServerServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

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
