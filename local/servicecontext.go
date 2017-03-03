package local

/**
 * Alter a service, based on the app.
 *
 * This primarily offers a way to use short forms in
 * service definitions in yml files, but is is primarily
 * targets at operations not based on compose, as compose
 * has expected behaviour, and is piped right through
 * the libCompose code.
 */

import (
	"strings"

	docker_cli_compose_types "github.com/docker/docker/cli/compose/types"

	handler_dockercli_stack "github.com/wunderkraut/radi-handler-dockercli/stack"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

// ServiceAlter alters a service based on internal conditions (for example rewriting volumes)
type LocalServiceContext struct {
	localSettings *handler_local.LocalAPISettings
}

// Constructor for LocalServiceContext
func New_LocalServiceContext(settings *handler_local.LocalAPISettings) *LocalServiceContext {
	return &LocalServiceContext{
		localSettings: settings,
	}
}

// Convert this to a stack ServiceAlter interface
func (localContext *LocalServiceContext) ServiceContext() handler_dockercli_stack.ServiceContext {
	return handler_dockercli_stack.ServiceContext(localContext)
}

// Provide a workingDir for compose relational mapping
func (localContext *LocalServiceContext) WorkingDir() string {
	return localContext.localSettings.ProjectRootPath
}

// Provide an ENV map for compose mapping
func (localContext *LocalServiceContext) EnvMap() map[string]string {
	return map[string]string{}
}

/**
 * clean up a service based on this app
 */
func (localContext *LocalServiceContext) AlterService(service *docker_cli_compose_types.ServiceConfig) {
	localContext.alterService_ProjectNetwork(service)
	localContext.alterService_RewriteMappedVolumes(service)
}

// make sure that a service is using the default network [@TODO THIS SHOULD NOT BE NECESSARY]
func (localContext *LocalServiceContext) alterService_ProjectNetwork(service *docker_cli_compose_types.ServiceConfig) {
	/**
	 * If a service has no network then we create the default network config.
	 *
	 * This is copypasta from github.com/docker/libcompose/project::Project.handleNetworkConfig()
	 * which means that we are duplicating internal functionality that may not be stable.
	 *
	 * This requirement came up after an update to the libcompose upstream library, which broke
	 * the existing missing network setup.  What is happening is that we are alteting a serviceconfig
	 * which will be added to a libcompose.Project::Project struct, and that struct has already
	 * run its initializer, which does the default network configuration.  This means that it is
	 * too late to simply add an empty network.  An alternative is to re-run the initializer, but
	 * as we have no access to that functionality from the Interface, there is little we can do.
	 */

	// if service.Networks == nil || len(service.Networks.Networks) == 0 {
	// 	// Add default as network
	// 	service.Networks = &libCompose_yaml.Networks{
	// 		Networks: []*libCompose_yaml.Network{
	// 			{
	// 				Name:     "default",
	// 				RealName: project.composeContext.Context.ProjectName + "_default",
	// 			},
	// 		},
	// 	}
	// }
}

// rewrite mapped service volumes to use app points.
func (localContext *LocalServiceContext) alterService_RewriteMappedVolumes(service *docker_cli_compose_types.ServiceConfig) {

	for index, volumeString := range service.Volumes {
		volumeParts := strings.Split(volumeString, ":")
		volumeSource := string(volumeParts[0])
		volumeSourceSplit := strings.SplitAfterN(volumeSource, "", 2) // we use this to look at the first character

		switch volumeSourceSplit[0] { // switch on the first char as a string

		// relate volume to the current user home path
		case "~":
			homePath := localContext.localSettings.UserHomePath
			volumeSource = strings.Replace(volumeSource, "~", homePath, 1)

		// relate volume to project root
		case ".":
			appPath := localContext.localSettings.ProjectRootPath
			volumeSource = strings.Replace(volumeSource, "~", appPath, 1)

		// @TODO this is a stupid special hard-code that we should document somehow
		// @NOTE this is dangerous and will likely only work in cases where PWD is available
		case "!":
			appPath := localContext.localSettings.ExecPath
			volumeSource = strings.Replace(volumeSource, "!", appPath, 1)

		case "@":
			if aliasPath, found := localContext.localSettings.ConfigPaths.Get(volumeSourceSplit[1]); found {
				volumeSource = strings.Replace(volumeSource, volumeSource, aliasPath.PathString(), 1)
			}

		default:
			continue
		}

		volumeParts[0] = volumeSource
		service.Volumes[index] = strings.Join(volumeParts, ":")
	}

}
