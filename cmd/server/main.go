package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juliendoutre/recorder/internal/server"
	v1 "github.com/juliendoutre/recorder/pkg/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	GoVersion string
	Os        string
	Arch      string
)

func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}
	defer func() { _ = logger.Sync() }()

	creds, err := credentials.NewServerTLSFromFile("/etc/recorder/server.crt.pem", "/etc/recorder/server.key.pem")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	}

	grpcServer := grpc.NewServer(grpcOptions...)

	healthgrpc.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)

	pgQuery := url.Values{}
	pgQuery.Add("sslmode", "verify-full")

	pgURL := url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		User:     url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Path:     os.Getenv("POSTGRES_DB"),
		RawQuery: pgQuery.Encode(),
	}

	pg, err := pgxpool.New(context.Background(), pgURL.String())
	if err != nil {
		logger.Panic("Connecting to DB", zap.Error(err))
	}
	defer pg.Close()

	v1.RegisterRecorderServer(grpcServer, server.New(&v1.Version{
		GoVersion: GoVersion,
		Os:        Os,
		Arch:      Arch,
	}, pg))

	grpcPort, err := strconv.ParseInt(os.Getenv("GRPC_PORT"), 10, 64)
	if err != nil {
		logger.Panic("Parsing gRPC port", zap.Error(err))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Panic("Failed to create a TCP listener", zap.Error(err))
	}

	go handleSignals(logger, grpcServer)

	logger.Info("Starting the Recorder server...", zap.Int64("port", grpcPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Panic("Failed to serve gRPC request", zap.Error(err))
	}
}

func handleSignals(logger *zap.Logger, server *grpc.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for signal := range signals {
		logger.Warn("Caught a cancellation signal, terminating...", zap.String("signal", signal.String()))
		server.GracefulStop()

		return
	}
}
