package stack

import (
	log "github.com/Sirupsen/logrus"

	api_property "github.com/wunderkraut/radi-api/property"
	api_usage "github.com/wunderkraut/radi-api/usage"

	handler_dockercli_stack_imported "github.com/wunderkraut/radi-handler-dockercli/stack/stack"
)

const (
	OPERATION_PROPERTY_DOCKER_STACK_DEPLOYOPTIONS_KEY = "docker.cli.command.stack.deployoptions"
	OPERATION_PROPERTY_DOCKER_STACK_REMOVEOPTIONS_KEY = "docker.cli.command.stack.removeoptions"
)

type DockercliStackDeployOptionsProperty struct {
	value handler_dockercli_stack_imported.DeployOptions
}

// Id for the property
func (depOpts *DockercliStackDeployOptionsProperty) Id() string {
	return OPERATION_PROPERTY_DOCKER_STACK_DEPLOYOPTIONS_KEY
}

// Id for the property
func (depOpts *DockercliStackDeployOptionsProperty) Type() string {
	return "github.com/docker/docker/cli/commands/stack.deployOptions"
}

// Label for the property
func (depOpts *DockercliStackDeployOptionsProperty) Label() string {
	return "Docker:Stack: Deploy options."
}

// Description for the property
func (depOpts *DockercliStackDeployOptionsProperty) Description() string {
	return "Deploy options for a docker stack command"
}

// Is the Property internal only
func (depOpts *DockercliStackDeployOptionsProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

// Property accessors
func (depOpts *DockercliStackDeployOptionsProperty) Get() interface{} {
	return interface{}(depOpts.value)
}
func (depOpts *DockercliStackDeployOptionsProperty) Set(value interface{}) bool {
	if converted, ok := value.(handler_dockercli_stack_imported.DeployOptions); ok {
		depOpts.value = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected github.com/wunderkraut/radi-handler-dockercli/stack/stack.DeployOptions struct")
		return false
	}
}

// Copy the property
func (depOpts *DockercliStackDeployOptionsProperty) Copy() api_property.Property {
	prop := &DockercliStackDeployOptionsProperty{}
	prop.Set(depOpts.Get())
	return api_property.Property(prop)
}

type DockercliStackRemoveOptionsProperty struct {
	value handler_dockercli_stack_imported.RemoveOptions
}

// Id for the property
func (remOpts *DockercliStackRemoveOptionsProperty) Id() string {
	return OPERATION_PROPERTY_DOCKER_STACK_REMOVEOPTIONS_KEY
}

// Id for the property
func (remOpts *DockercliStackRemoveOptionsProperty) Type() string {
	return "github.com/docker/docker/cli/commands/stack.deployOptions"
}

// Label for the property
func (remOpts *DockercliStackRemoveOptionsProperty) Label() string {
	return "Docker:Stack: Deploy options."
}

// Description for the property
func (remOpts *DockercliStackRemoveOptionsProperty) Description() string {
	return "Deploy options for a docker stack command"
}

// Is the Property internal only
func (remOpts *DockercliStackRemoveOptionsProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

// Property accessors
func (remOpts *DockercliStackRemoveOptionsProperty) Get() interface{} {
	return interface{}(remOpts.value)
}
func (remOpts *DockercliStackRemoveOptionsProperty) Set(value interface{}) bool {
	if converted, ok := value.(handler_dockercli_stack_imported.RemoveOptions); ok {
		remOpts.value = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected github.com/wunderkraut/radi-handler-dockercli/stack/stack.RemoveOptions struct")
		return false
	}
}

// Copy the property
func (remOpts *DockercliStackRemoveOptionsProperty) Copy() api_property.Property {
	prop := &DockercliStackRemoveOptionsProperty{}
	prop.Set(remOpts.Get())
	return api_property.Property(prop)
}
