package local

type DockercliLocalHandlerBase struct {
	config DockercliLocalConfig
}

func New_DockercliLocalHandlerBase(config DockercliLocalConfig) *DockercliLocalHandlerBase {
	return &DockercliLocalHandlerBase{
		config: config,
	}
}

func (base *DockercliLocalHandlerBase) DockercliLocalConfig() DockercliLocalConfig {
	return base.config
}

func (base *DockercliLocalHandlerBase) SetDockercliLocalConfig(config DockercliLocalConfig) {
	base.config = config
}
