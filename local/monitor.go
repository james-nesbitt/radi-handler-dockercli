package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_result "github.com/wunderkraut/radi-api/result"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
	handler_dockercli_stack "github.com/wunderkraut/radi-handler-dockercli/stack"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

/**
 * Some dockercli implementation monitor methods for things like outputting
 * settings for testing
 */

// Local Handler for monitoring using docker cli
type MonitorHandler struct {
	handler_local.LocalHandler_Base
	DockercliLocalHandlerBase
	handler_dockercli.DockercliHandlerBase
	handler_dockercli_stack.DockercliStackHandlerBase
}

// Constructor for MonitorHandler
func New_MonitorHandler(localBase *handler_local.LocalHandler_Base, dockerCLIBase *handler_dockercli.DockercliHandlerBase, stackBase *handler_dockercli_stack.DockercliStackHandlerBase) *MonitorHandler {
	return &MonitorHandler{
		LocalHandler_Base:         *localBase,
		DockercliHandlerBase:      *dockerCLIBase,
		DockercliStackHandlerBase: *stackBase,
	}
}

// Validate the Base Handler
func (base *MonitorHandler) Id() string {
	return "dockercli.monitor"
}

// Validate the Base Handler
func (base *MonitorHandler) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Validate the Base Handler
func (base *MonitorHandler) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// use a single base operation
	baseCliOp := base.DockercliOperationBase()
	baseStackOp := base.DockercliStackOperationBase()

	ops.Add(api_operation.Operation(&handler_dockercli_stack.DockercliStackMonitorPsOperation{
		DockercliOperationBase:      *baseCliOp,
		DockercliStackOperationBase: *baseStackOp,
	}))

	return ops.Operations()
}
