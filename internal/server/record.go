package server

import (
	"context"
	"encoding/json"

	v1 "github.com/juliendoutre/recorder/pkg/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Record(ctx context.Context, input *v1.RecordInput) (*emptypb.Empty, error) {
	token, err := s.jwtParser.Parse(input.GetJwt(), s.jwkStore.Keyfunc)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "TODO")
	}

	if !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "TODO")
	}

	headers, err := json.Marshal(token.Header)
	if err != nil {
		return nil, status.Error(codes.Internal, "TODO")
	}

	payload, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, status.Error(codes.Internal, "TODO")
	}

	if _, err := s.pg.Exec(ctx, "INSERT INTO recorder.claims (header, payload) VALUES ($1, $2);", headers, payload); err != nil {
		return nil, status.Error(codes.Internal, "TODO")
	}

	return &emptypb.Empty{}, nil
}
