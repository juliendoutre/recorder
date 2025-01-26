package server

import (
	"context"

	v1 "github.com/juliendoutre/recorder/pkg/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Record(ctx context.Context, input *v1.RecordInput) (*emptypb.Empty, error) {
	token, err := s.jwtParser.Parse(input.GetJwt(), s.jwkStore.Keyfunc)
	if err != nil {
		// TODO: return error
	}

	if !token.Valid {
		// TODO: return error
	}

	// TODO: insert to DB

	return &emptypb.Empty{}, nil
}
