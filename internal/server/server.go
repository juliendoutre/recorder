package server

import (
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/juliendoutre/recorder/pkg/v1"
)

func New(version *v1.Version, pg *pgxpool.Pool, jwkStore keyfunc.Keyfunc) (*Server, error) {
	return &Server{
		version:   version,
		pg:        pg,
		jwtParser: jwt.NewParser(),
		jwkStore:  jwkStore,
	}, nil
}

type Server struct {
	v1.UnimplementedRecorderServer

	version   *v1.Version
	pg        *pgxpool.Pool
	jwtParser *jwt.Parser
	jwkStore  keyfunc.Keyfunc
}
