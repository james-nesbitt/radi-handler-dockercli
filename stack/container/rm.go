package container

import (
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type RmOptions struct {
	rmVolumes bool
	rmLink    bool
	force     bool

	containers []string
}

func runRm(dockerCli *command.DockerCli, opts *RmOptions) error {
	ctx := context.Background()

	var errs []string
	options := types.ContainerRemoveOptions{
		RemoveVolumes: opts.rmVolumes,
		RemoveLinks:   opts.rmLink,
		Force:         opts.force,
	}

	errChan := parallelOperation(ctx, opts.containers, func(ctx context.Context, container string) error {
		container = strings.Trim(container, "/")
		if container == "" {
			return errors.New("Container name cannot be empty")
		}
		return dockerCli.Client().ContainerRemove(ctx, container, options)
	})

	for _, name := range opts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintln(dockerCli.Out(), name)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
