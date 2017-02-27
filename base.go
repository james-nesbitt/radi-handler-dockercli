package dockercli

import (
	"io"

	docker_cli_command "github.com/docker/docker/cli/command"
	docker_cli_flags "github.com/docker/docker/cli/flags"

	api_result "github.com/wunderkraut/radi-api/result"
)

/**
 * Base handlers and operations
 */

type DockercliHandlerBase struct {
	in   io.ReadCloser
	out  io.Writer
	err  io.Writer
	opts *docker_cli_flags.ClientOptions

	cli *docker_cli_command.DockerCli

	DockercliOperationBase_common *DockercliOperationBase
}

// Constructor for DockercliOperationBase
func New_DockercliHandlerBase(in io.ReadCloser, out, err io.Writer, opts *docker_cli_flags.ClientOptions) *DockercliHandlerBase {
	return &DockercliHandlerBase{
		in:   in,
		out:  out,
		err:  err,
		opts: opts,
	}
}

// Validate the Base Handler
func (base *DockercliHandlerBase) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Retreive the Docker Cli
func (base *DockercliHandlerBase) DockerCli() *docker_cli_command.DockerCli {
	if base.cli == nil {
		base.cli = docker_cli_command.NewDockerCli(base.in, base.out, base.err)
		base.cli.Initialize(base.opts)
	}
	return base.cli
}

// Retreive the Docker Cli
func (base *DockercliHandlerBase) DockercliOperationBase() *DockercliOperationBase {
	if base.DockercliOperationBase_common == nil {
		base.DockercliOperationBase_common = New_DockercliOperationBase(base.DockerCli())
	}
	return base.DockercliOperationBase_common
}

/**
 * Base CLI operation with a CLI generator.
 */

type DockercliOperationBase struct {
	cli *docker_cli_command.DockerCli
}

// Constructor for DockercliOperationBase
func New_DockercliOperationBase(cli *docker_cli_command.DockerCli) *DockercliOperationBase {
	return &DockercliOperationBase{
		cli: cli,
	}
}

// Retreive the Docker Cli
func (base *DockercliOperationBase) DockerCli() *docker_cli_command.DockerCli {
	return base.cli
}

// Assign the Docker Cli
func (base *DockercliOperationBase) SetCli(cli *docker_cli_command.DockerCli) {
	base.cli = cli
}
