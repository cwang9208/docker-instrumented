package container

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type rmOptions struct {
	rmVolumes bool
	rmLink    bool
	force     bool

	containers []string
}

// NewRmCommand creates a new cobra.Command for `docker rm`
func NewRmCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts rmOptions

	cmd := &cobra.Command{
		Use:   "rm [OPTIONS] CONTAINER [CONTAINER...]",
		Short: "Remove one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runRm(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.rmVolumes, "volumes", "v", false, "Remove the volumes associated with the container")
	flags.BoolVarP(&opts.rmLink, "link", "l", false, "Remove the specified link")
	flags.BoolVarP(&opts.force, "force", "f", false, "Force the removal of a running container (uses SIGKILL)")
	return cmd
}

func runRm(dockerCli *command.DockerCli, opts *rmOptions) error {
	ctx := context.Background()
	stderr := dockerCli.Err()

	var (
		errs []string
		timeStart time.Time
		timeDone  time.Time
	)

	timeStart = time.Now()

	options := types.ContainerRemoveOptions{
		RemoveVolumes: opts.rmVolumes,
		RemoveLinks:   opts.rmLink,
		Force:         opts.force,
	}

	errChan := parallelOperation(ctx, opts.containers, func(ctx context.Context, container string) error {
		if container == "" {
			return fmt.Errorf("Container name cannot be empty")
		}
		container = strings.Trim(container, "/")
		return dockerCli.Client().ContainerRemove(ctx, container, options)
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
	timeDone = time.Now()
	fmt.Fprintf(stderr, "RM Start: %s\n", timeStart.Format(time.RFC3339Nano))
	fmt.Fprintf(stderr, "RM Done:  %s\n", timeDone.Format(time.RFC3339Nano))
	fmt.Fprintf(stderr, "(note this is only accurate with single-container rm operations)\n")
	fmt.Fprintf(stderr, "Duration:  %d nanoseconds\n", timeDone.Sub(timeStart).Nanoseconds())

	return nil
}
