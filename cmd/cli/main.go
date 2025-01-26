package main

import (
	"os"

	"github.com/juliendoutre/recorder/internal/cli"
	v1 "github.com/juliendoutre/recorder/pkg/v1"
)

var (
	GoVersion string
	Os        string
	Arch      string
)

func main() {
	if err := cli.RootCmd(&v1.Version{
		GoVersion: GoVersion,
		Os:        Os,
		Arch:      Arch,
	}).Execute(); err != nil {
		os.Exit(1)
	}
}
