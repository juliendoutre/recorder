package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/juliendoutre/recorder/pkg/v1"
)

func New(version *v1.Version, pg *pgxpool.Pool) *Server {
	return &Server{version: version, pg: pg}
}

type Server struct {
	v1.UnimplementedRecorderServer

	version *v1.Version
	pg      *pgxpool.Pool
}
