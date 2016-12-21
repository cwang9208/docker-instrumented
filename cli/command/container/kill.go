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

type killOptions struct {
	signal string

	containers []string
}

// NewKillCommand creates a new cobra.Command for `docker kill`
func NewKillCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts killOptions

	cmd := &cobra.Command{
		Use:   "kill [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Kill one or more running containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runKill(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.signal, "signal", "s", "KILL", "Signal to send to the container")
	return cmd
}

func runKill(dockerCli *command.DockerCli, opts *killOptions) error {
	var errs []string
	var (
		timeStart time.Time
		timeDone  time.Time
	)

	timeStart = time.Now()

	ctx := context.Background()
	errChan := parallelOperation(ctx, opts.containers, func(ctx context.Context, container string) error {
		retval := dockerCli.Client().ContainerKill(ctx, container, opts.signal)
		timeDone = time.Now()
		fmt.Fprintf(dockerCli.Err(), "Kill Start: %s\n", timeStart.Format(time.RFC3339Nano))
		fmt.Fprintf(dockerCli.Err(), "Kill Done:  %s\n", timeDone.Format(time.RFC3339Nano))
		fmt.Fprintf(dockerCli.Err(), "Duration:  %d nanoseconds\n", timeDone.Sub(timeStart).Nanoseconds())
		return retval
	})
	for _, name := range opts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
		} else {
			fmt.Fprintf(dockerCli.Out(), "%s\n", name)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}
