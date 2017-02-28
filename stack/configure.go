package stack

import (
	handler_dockercli_stack_imported "github.com/wunderkraut/radi-handler-dockercli/stack/stack" // "github.com/docker/docker/cli/command/stack"
)

/**
 * A configuring interface for the stack handlers
 */

type DockercliStackConfig interface {
	DeployOptions() *handler_dockercli_stack_imported.DeployOptions
	RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions
	PsOptions() *handler_dockercli_stack_imported.PsOptions
}
