package server

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	v1 "github.com/juliendoutre/recorder/pkg/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var digestRegex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

func (s *Server) Record(ctx context.Context, input *v1.RecordInput) (*emptypb.Empty, error) {
	digest := input.GetDigest()

	if !digestRegex.MatchString(input.GetDigest()) {
		return nil, status.Error(codes.InvalidArgument, "invalid digest") //nolint:wrapcheck
	}

	token, err := s.jwtParser.Parse(input.GetJwt(), s.jwkStore.Keyfunc)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("parsing JWT: %s", err)) //nolint:wrapcheck
	}

	if !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token") //nolint:wrapcheck
	}

	headers, err := json.Marshal(token.Header)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("marshalling headers: %s", err)) //nolint:wrapcheck
	}

	payload, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("marshalling payload: %s", err)) //nolint:wrapcheck
	}

	if _, err := s.pg.Exec(
		ctx,
		"INSERT INTO recorder.claims (digest, header, payload) VALUES ($1, $2, $3);",
		digest, headers, payload,
	); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("inserting record: %s", err)) //nolint:wrapcheck
	}

	return &emptypb.Empty{}, nil
}
