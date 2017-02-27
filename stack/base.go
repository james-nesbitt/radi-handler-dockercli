package stack

/**
 * Base operation for all stack operations
 */

type DockercliStackHandlerBase struct {
	configure DockercliStackConfig
}

// Constructor for DockercliStackOperationBase
func New_DockercliStackHandlerBase(configure DockercliStackConfig) *DockercliStackHandlerBase {
	return &DockercliStackHandlerBase{
		configure: configure,
	}
}

func (stackBase *DockercliStackHandlerBase) SetDockercliStackConfig(configure DockercliStackConfig) {
	stackBase.configure = configure
}

func (stackBase *DockercliStackHandlerBase) DockercliStackConfig() DockercliStackConfig {
	return stackBase.configure
}

func (stackBase *DockercliStackHandlerBase) DockercliStackOperationBase() *DockercliStackOperationBase {
	return New_DockercliStackOperationBase(stackBase.DockercliStackConfig())
}

/**
 * Base operation for all stack operations
 */

type DockercliStackOperationBase struct {
	configure DockercliStackConfig
}

// Constructor for DockercliStackOperationBase
func New_DockercliStackOperationBase(configure DockercliStackConfig) *DockercliStackOperationBase {
	return &DockercliStackOperationBase{
		configure: configure,
	}
}

func (stackBase *DockercliStackOperationBase) SetDockercliStackConfig(configure DockercliStackConfig) {
	stackBase.configure = configure
}

func (stackBase *DockercliStackOperationBase) DockercliStackConfig() DockercliStackConfig {
	return stackBase.configure
}

func (stackBase *DockercliStackOperationBase) DeployOptionsProperty() *DockercliStackDeployOptionsProperty {
	deployOpts := stackBase.DockercliStackConfig().DeployOptions()
	deployOptsProp := DockercliStackDeployOptionsProperty{}
	deployOptsProp.Set(*deployOpts)
	return &deployOptsProp
}

func (stackBase *DockercliStackOperationBase) RemoveOptionsProperty() *DockercliStackRemoveOptionsProperty {
	remOpts := stackBase.DockercliStackConfig().RemoveOptions()
	remOptsProp := DockercliStackRemoveOptionsProperty{}
	remOptsProp.Set(*remOpts)
	return &remOptsProp
}
