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
type DockercliOrchestrateHandler struct {
	handler_local.LocalHandler_Base
	DockercliLocalHandlerBase
	handler_dockercli.DockercliHandlerBase
	handler_dockercli_stack.DockercliStackHandlerBase
}

// Constructor for DockercliOrchestrateHandler
func New_DockercliOrchestrateHandler(localBase *handler_local.LocalHandler_Base, dockerCLIBase *handler_dockercli.DockercliHandlerBase, stackBase *handler_dockercli_stack.DockercliStackHandlerBase) *DockercliOrchestrateHandler {
	return &DockercliOrchestrateHandler{
		LocalHandler_Base:         *localBase,
		DockercliHandlerBase:      *dockerCLIBase,
		DockercliStackHandlerBase: *stackBase,
	}
}

// Validate the Base Handler
func (base *DockercliOrchestrateHandler) Id() string {
	return "dockercli.orchestrate"
}

// Validate the Base Handler
func (base *DockercliOrchestrateHandler) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Validate the Base Handler
func (base *DockercliOrchestrateHandler) Operations() api_operation.Operations {
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
