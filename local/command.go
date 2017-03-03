package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_result "github.com/wunderkraut/radi-api/result"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
	handler_dockercli_stack "github.com/wunderkraut/radi-handler-dockercli/stack"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

/**
 * Some dockercli implementation command methods for things like outputting
 * settings for testing
 */

// Local Handler for commanding using docker cli
type CommandHandler struct {
	handler_local.LocalHandler_Base
	handler_local.LocalHandler_ConfigWrapperBase
	DockercliLocalHandlerBase
	handler_dockercli.DockercliHandlerBase
	handler_dockercli_stack.DockercliStackHandlerBase
}

// Constructor for CommandHandler
func New_CommandHandler(localBase *handler_local.LocalHandler_Base, dockerCLIBase *handler_dockercli.DockercliHandlerBase, stackBase *handler_dockercli_stack.DockercliStackHandlerBase) *CommandHandler {
	return &CommandHandler{
		LocalHandler_Base:         *localBase,
		DockercliHandlerBase:      *dockerCLIBase,
		DockercliStackHandlerBase: *stackBase,
	}
}

// Validate the Base Handler
func (base *CommandHandler) Id() string {
	return "dockercli.command"
}

// Validate the Base Handler
func (base *CommandHandler) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Validate the Base Handler
func (base *CommandHandler) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// use shared base operations
	baseCliOp := base.DockercliOperationBase()
	baseStackOp := base.DockercliStackOperationBase()

	// build the command config wrapper (which will pull from config yaml)
	configWrapper := base.ConfigWrapper()    // from LocalHandler_Base
	localSettings := base.LocalAPISettings() // from LocalHandler_Base
	localServiceContext := New_LocalServiceContext(localSettings).ServiceContext()
	commandConfigConnector := handler_dockercli_stack.New_DockercliStackCommand_ConfigureYml(configWrapper, localServiceContext, baseCliOp, baseStackOp).CommandConfigConnector()

	ops.Add(api_operation.Operation(&handler_dockercli_stack.DockercliStackCommandListOperation{
		DockercliOperationBase:      *baseCliOp,
		DockercliStackOperationBase: *baseStackOp,
		Connector:                   commandConfigConnector,
	}))
	ops.Add(api_operation.Operation(&handler_dockercli_stack.DockercliStackCommandGetOperation{
		DockercliOperationBase:      *baseCliOp,
		DockercliStackOperationBase: *baseStackOp,
		Connector:                   commandConfigConnector,
	}))

	return ops.Operations()
}
