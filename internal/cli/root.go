package cli

import (
	"fmt"
	"time"

	v1 "github.com/juliendoutre/recorder/pkg/v1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	host string
	port int

	conn   *grpc.ClientConn
	client v1.RecorderClient
	logger *zap.Logger
)

func RootCmd(version *v1.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "recorder",
		Short:        "A CLI command to interact with Recorder.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			logger, err = zap.NewProductionConfig().Build()
			if err != nil {
				return err
			}

			clientOptions := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithConnectParams(grpc.ConnectParams{
					MinConnectTimeout: 1 * time.Second,
					Backoff:           backoff.DefaultConfig,
				}),
			}

			conn, err = grpc.NewClient(fmt.Sprintf("%s:%d", host, port), clientOptions...)
			if err != nil {
				return err
			}

			client = v1.NewRecorderClient(conn)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if err := conn.Close(); err != nil {
				return err
			}

			_ = logger.Sync()

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&host, "host", "localhost", "Host the Recorder server listens on.")
	cmd.PersistentFlags().IntVar(&port, "port", 8000, "Port the Recorder server listens on.")

	cmd.AddCommand(
		versionCmd(version),
	)

	return cmd
}
