package cli

import (
	"encoding/json"
	"fmt"

	v1 "github.com/juliendoutre/recorder/pkg/v1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func versionCmd(v *v1.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print out the CLI version.",
		RunE: func(cmd *cobra.Command, args []string) error {
			versions := map[string]*v1.Version{
				"client": v,
			}

			serverVersion, err := client.GetVersion(cmd.Context(), &emptypb.Empty{})
			if err != nil {
				logger.Debug("Could not reach out to the server", zap.Error(err))
			} else {
				versions["server"] = serverVersion
			}

			data, err := json.MarshalIndent(versions, "", "    ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))

			return nil
		},
	}

	return cmd
}
