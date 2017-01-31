package container

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type pauseOptions struct {
	containers []string
}

// NewPauseCommand creates a new cobra.Command for `docker pause`
func NewPauseCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts pauseOptions

	return &cobra.Command{
		Use:   "pause CONTAINER [CONTAINER...]",
		Short: "Pause all processes within one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runPause(dockerCli, &opts)
		},
	}
}

func runPause(dockerCli *command.DockerCli, opts *pauseOptions) error {
	var (
		timeStart time.Time
		timeDone  time.Time
	)

	timeStart = time.Now()
	stderr := dockerCli.Err()

	ctx := context.Background()

	var errs []string
	errChan := parallelOperation(ctx, opts.containers, dockerCli.Client().ContainerPause)
	for _, container := range opts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
		} else {
			fmt.Fprintf(dockerCli.Out(), "%s\n", container)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	timeDone = time.Now()
	fmt.Fprintf(stderr, "pause Start: %s\n", timeStart.Format(time.RFC3339Nano))
	fmt.Fprintf(stderr, "pause Done:  %s\n", timeDone.Format(time.RFC3339Nano))
	fmt.Fprintf(stderr, "(note this is only accurate with single-container rm operations)\n")
	fmt.Fprintf(stderr, "Duration:  %d nanoseconds\n", timeDone.Sub(timeStart).Nanoseconds())

	return nil
}
