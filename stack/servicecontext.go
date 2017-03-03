package stack

/**
 * Alter a service, based on the app.
 */

import (
	docker_cli_compose_types "github.com/docker/docker/cli/compose/types"
)

// ServiceAlter provides provides contextual service functionality on internal conditions (for example rewriting volumes)
type ServiceContext interface {
	// Provide a workingDir for compose relational mapping
	WorkingDir() string
	// Provide an ENV map for compose mapping
	EnvMap() map[string]string
	//clean up a service based on this app
	AlterService(service *docker_cli_compose_types.ServiceConfig)
}
