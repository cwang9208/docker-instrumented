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

type unpauseOptions struct {
	containers []string
}

// NewUnpauseCommand creates a new cobra.Command for `docker unpause`
func NewUnpauseCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts unpauseOptions

	cmd := &cobra.Command{
		Use:   "unpause CONTAINER [CONTAINER...]",
		Short: "Unpause all processes within one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runUnpause(dockerCli, &opts)
		},
	}
	return cmd
}

func runUnpause(dockerCli *command.DockerCli, opts *unpauseOptions) error {
	var (
		timeStart time.Time
		timeDone  time.Time
	)

	timeStart = time.Now()
	stderr := dockerCli.Err()

	ctx := context.Background()

	var errs []string
	errChan := parallelOperation(ctx, opts.containers, dockerCli.Client().ContainerUnpause)
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
	fmt.Fprintf(stderr, "unpause Start: %s\n", timeStart.Format(time.RFC3339Nano))
	fmt.Fprintf(stderr, "unpause Done:  %s\n", timeDone.Format(time.RFC3339Nano))
	fmt.Fprintf(stderr, "(note this is only accurate with single-container rm operations)\n")
	fmt.Fprintf(stderr, "Duration:  %d nanoseconds\n", timeDone.Sub(timeStart).Nanoseconds())

	return nil
}
