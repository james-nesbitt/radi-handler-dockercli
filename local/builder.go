package local

import (
	log "github.com/Sirupsen/logrus"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_result "github.com/wunderkraut/radi-api/result"

	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"

	handler_local "github.com/wunderkraut/radi-handlers/local"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
	handler_dockercli_stack "github.com/wunderkraut/radi-handler-dockercli/stack"
)

/**
 * A Local build implementation for providing local operations using the
 * docker cli various implementations
 */

// The api Builder for this dockercli local
type LocalBuilder struct {
	handler_local.LocalBuilder
	settings handler_local.LocalAPISettings

	dockercliConfig DockercliLocalConfig

	parent api_api.API

	DockercliHandlerBase_common      *handler_dockercli.DockercliHandlerBase
	DockercliStackHandlerBase_common *handler_dockercli_stack.DockercliStackHandlerBase
}

// Constructor for LocalBuilder
func New_LocalBuilder(settings handler_local.LocalAPISettings, dockercliConfig DockercliLocalConfig) *LocalBuilder {
	return &LocalBuilder{
		LocalBuilder:    *handler_local.New_LocalBuilder(settings),
		settings:        settings,
		dockercliConfig: dockercliConfig,
	}
}

// IBuilder ID
func (builder *LocalBuilder) Id() string {
	return "dockercli_local"
}

// Builder Settings
func (builder *LocalBuilder) LocalAPISettings() handler_local.LocalAPISettings {
	return builder.settings
}

// Set the parent API, which may need to build Config and Setting Wrappers
func (builder *LocalBuilder) SetAPI(parent api_api.API) {
	builder.parent = parent
	builder.LocalBuilder.SetAPI(parent)
}

// Initialize the handler for certain implementations
func (builder *LocalBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) api_result.Result {
	dockercliConfig := builder.DockercliConfig(settingsProvider)

	localBase := builder.Base()
	dockerCLIBase := builder.base_DockercliHandlerBase(dockercliConfig)
	stackBase := builder.base_DockercliStackHandlerBase(dockercliConfig)

	/**
	 * Here you could override that bases depending on the settings
	 */

	for _, implementation := range implementations.Order() {
		switch implementation {
		case "orchestrate":
			builder.build_Orchestrate(localBase, dockerCLIBase, stackBase)
		default:
			log.WithFields(log.Fields{"implementation": implementation}).Warn("Local builder implementation not available")
		}
	}

	return api_result.MakeSuccessfulResult()
}

/**
 * Common creator of configuration tools used for various implementations
 */

// Add local Handlers for Config and Settings
func (builder *LocalBuilder) ConfigWrapper() api_config.ConfigWrapper {
	// Build a configWrapper if needed
	if builder.Config == nil {
		builder.Config = api_config.New_SimpleConfigWrapper(builder.parent.Operations())
	}
	return builder.Config
}

// Add local Handlers for Setting
func (builder *LocalBuilder) SettingWrapper() api_setting.SettingWrapper {
	// Build a configWrapper if needed
	if builder.Setting == nil {
		builder.Setting = api_setting.New_SimpleSettingWrapper(builder.parent.Operations())
	}
	return builder.Setting
}

// Build a cli config for this builder (currently just builds a single yml based configure, regardless of any settings.)
func (builder *LocalBuilder) DockercliConfig(settingsProvider api_builder.SettingsProvider) DockercliLocalConfig {
	if builder.dockercliConfig == nil {
		settings := builder.LocalAPISettings()
		settingWrapper := builder.SettingWrapper()
		// configWrapper := builder.ConfigWrapper()

		builder.dockercliConfig = New_DockercliLocalConfigDefault(settings, settingWrapper)
		// builder.dockercliConfig := New_DockercliLocalConfigConfigWrapperYml(settings, configWrapper)
	}
	return builder.dockercliConfig
}

/**
 * Common base builders
 */

func (builder *LocalBuilder) base_DockercliHandlerBase(dockercliConfig DockercliLocalConfig) *handler_dockercli.DockercliHandlerBase {
	if builder.DockercliHandlerBase_common == nil {
		in, out, err := dockercliConfig.IO()
		opts := dockercliConfig.ClientOptions()

		builder.DockercliHandlerBase_common = handler_dockercli.New_DockercliHandlerBase(in, out, err, opts)
	}
	return builder.DockercliHandlerBase_common
}

func (builder *LocalBuilder) base_DockercliStackHandlerBase(dockercliConfig DockercliLocalConfig) *handler_dockercli_stack.DockercliStackHandlerBase {
	if builder.DockercliStackHandlerBase_common == nil {
		dockercliStackConfig := handler_dockercli_stack.DockercliStackConfig(builder.dockercliConfig)
		builder.DockercliStackHandlerBase_common = handler_dockercli_stack.New_DockercliStackHandlerBase(dockercliStackConfig)
	}
	return builder.DockercliStackHandlerBase_common
}

/**
 * Actual build abstractions per implemention
 */

// Build and add a handler for orchestration
func (builder *LocalBuilder) build_Orchestrate(localBase *handler_local.LocalHandler_Base, dockerCLIBase *handler_dockercli.DockercliHandlerBase, stackBase *handler_dockercli_stack.DockercliStackHandlerBase) api_result.Result {
	local_orchestration := New_DockercliOrchestrateHandler(localBase, dockerCLIBase, stackBase)

	res := local_orchestration.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(local_orchestration))
		// Get an orchestrate wrapper for other handlers
		//builder.Orchestrate = local_orchestration.OrchestrateWrapper()

		log.Debug("DockerCLI:localBuilder: Built Orchestrate handler")
	}

	return res
}
