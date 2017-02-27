package stack

import (
	log "github.com/Sirupsen/logrus"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"

	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
	handler_dockercli_stack_imported "github.com/wunderkraut/radi-handler-dockercli/stack/stack" // "github.com/docker/docker/cli/command/stack"
)

const (
	OPERATION_ID_DOCKERCLI_STACK_DOWN = "dockercli.stack.orchestrate.down"
)

/**
 * Orchestration operations
 */

// Base class for config list Operation
type DockercliStackOrchestrateDownOperation struct {
	api_orchestrate.BaseOrchestrationDownOperation
	handler_dockercli.DockercliOperationBase
	DockercliStackOperationBase
}

// Id the operation
func (down *DockercliStackOrchestrateDownOperation) Id() string {
	return OPERATION_ID_DOCKERCLI_STACK_DOWN
}

// Define the operations as externally used
func (down *DockercliStackOrchestrateDownOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Return Operation properties
func (down *DockercliStackOrchestrateDownOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	// Use a deploy Opts propperty, with a default set to the configured DeployOptis
	props.Add(api_property.Property(down.RemoveOptionsProperty()))

	return props.Properties()
}

// Validate the operation
func (down *DockercliStackOrchestrateDownOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Execute the operation
func (down *DockercliStackOrchestrateDownOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	go func() {
		optsProp, _ := props.Get(OPERATION_PROPERTY_DOCKER_STACK_REMOVEOPTIONS_KEY)
		opts := optsProp.Get().(handler_dockercli_stack_imported.RemoveOptions)

		cli := down.DockerCli()

		log.WithFields(log.Fields{"RemoveOptions": opts}).Info("Running Down orchestration using docker cli stack")

		if err := handler_dockercli_stack_imported.RunRemove(cli, opts); err == nil {
			res.MarkSuccess()
		} else {
			res.AddError(err)
			res.MarkFailed()
		}
		res.MarkFinished()
	}()

	return res.Result()
}
