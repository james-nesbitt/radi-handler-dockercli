package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_result "github.com/wunderkraut/radi-api/result"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
	handler_dockercli_stack "github.com/wunderkraut/radi-handler-dockercli/stack"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

/**
 * Handler for local orchestration through dockercli
 */

// Local Handler for orchestration using docker cli
type OrchestrateHandler struct {
	handler_local.LocalHandler_Base
	DockercliLocalHandlerBase
	handler_dockercli.DockercliHandlerBase
	handler_dockercli_stack.DockercliStackHandlerBase
}

// Constructor for OrchestrateHandler
func New_OrchestrateHandler(localBase *handler_local.LocalHandler_Base, dockerCLIBase *handler_dockercli.DockercliHandlerBase, stackBase *handler_dockercli_stack.DockercliStackHandlerBase) *OrchestrateHandler {
	return &OrchestrateHandler{
		LocalHandler_Base:         *localBase,
		DockercliHandlerBase:      *dockerCLIBase,
		DockercliStackHandlerBase: *stackBase,
	}
}

// Validate the Base Handler
func (base *OrchestrateHandler) Id() string {
	return "dockercli.orchestrate"
}

// Validate the Base Handler
func (base *OrchestrateHandler) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Validate the Base Handler
func (base *OrchestrateHandler) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// use a single base operation
	baseCliOp := base.DockercliOperationBase()
	baseStackOp := base.DockercliStackOperationBase()

	ops.Add(api_operation.Operation(&handler_dockercli_stack.DockercliStackOrchestrateUpOperation{
		DockercliOperationBase:      *baseCliOp,
		DockercliStackOperationBase: *baseStackOp,
	}))
	ops.Add(api_operation.Operation(&handler_dockercli_stack.DockercliStackOrchestrateDownOperation{
		DockercliOperationBase:      *baseCliOp,
		DockercliStackOperationBase: *baseStackOp,
	}))

	return ops.Operations()
}
