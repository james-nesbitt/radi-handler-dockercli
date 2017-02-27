package stack

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/cli/command"
)

const (
	defaultNetworkDriver = "overlay"
)

type DeployOptions struct {
	bundlefile       string
	composefile      string
	namespace        string
	sendRegistryAuth bool
}

func New_DeployOptions(bundlefile string, composefile string, namespace string, sendRegistryAuth bool) *DeployOptions {
	return &DeployOptions{
		bundlefile:       bundlefile,
		composefile:      composefile,
		namespace:        namespace,
		sendRegistryAuth: sendRegistryAuth,
	}
}

func RunDeploy(dockerCli *command.DockerCli, opts DeployOptions) error {
	ctx := context.Background()

	switch {
	case opts.bundlefile == "" && opts.composefile == "":
		return fmt.Errorf("Please specify either a bundle file (with --bundle-file) or a Compose file (with --compose-file).")
	case opts.bundlefile != "" && opts.composefile != "":
		return fmt.Errorf("You cannot specify both a bundle file and a Compose file.")
	case opts.bundlefile != "":
		return deployBundle(ctx, dockerCli, opts)
	default:
		return deployCompose(ctx, dockerCli, opts)
	}
}

// checkDaemonIsSwarmManager does an Info API call to verify that the daemon is
// a swarm manager. This is necessary because we must create networks before we
// create services, but the API call for creating a network does not return a
// proper status code when it can't create a network in the "global" scope.
func checkDaemonIsSwarmManager(ctx context.Context, dockerCli *command.DockerCli) error {
	info, err := dockerCli.Client().Info(ctx)
	if err != nil {
		return err
	}
	if !info.Swarm.ControlAvailable {
		return errors.New("This node is not a swarm manager. Use \"docker swarm init\" or \"docker swarm join\" to connect this node to swarm and try again.")
	}
	return nil
}
