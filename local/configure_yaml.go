package local

import (
	"errors"
	"io"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	api_config "github.com/wunderkraut/radi-api/operation/config"

	docker_cli_flags "github.com/docker/docker/cli/flags"

	handler_dockercli_stack_imported "github.com/wunderkraut/radi-handler-dockercli/stack/stack"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

const (
	// Config key used to retrieve docker cli local settings
	CONFIG_KEY_DOCKERCLI_LOCAL = "dockercli"
)

type DockercliLocalConfigConfigWrapperYml struct {
	DockercliLocalConfigDefault
	configWrapper api_config.ConfigWrapper

	config dockercliLocalConfigureYML
}

// Constructor for DockercliLocalConfigConfigWrapperYml
func New_DockercliLocalConfigConfigWrapperYml(localAPISettings handler_local.LocalAPISettings, configWrapper api_config.ConfigWrapper) *DockercliLocalConfigConfigWrapperYml {
	return &DockercliLocalConfigConfigWrapperYml{
		DockercliLocalConfigDefault: DockercliLocalConfigDefault{
			settings: localAPISettings,
		},
		configWrapper: configWrapper,
	}
}

// return this struct as a DockercliLocalConfig interface
func (configYml *DockercliLocalConfigConfigWrapperYml) DockercliLocalConfig() DockercliLocalConfig {
	return DockercliLocalConfig(configYml)
}

/**
 * DockercliLocalConfig Interface methods
 */

func (configYml *DockercliLocalConfigConfigWrapperYml) ClientOptions() *docker_cli_flags.ClientOptions {
	return docker_cli_flags.NewClientOptions()
}

func (configYml *DockercliLocalConfigConfigWrapperYml) DeployOptions() *handler_dockercli_stack_imported.DeployOptions {
	return handler_dockercli_stack_imported.New_DeployOptions(
		"",    // bundlefile,
		"",    // composefile,
		"",    // namespace,
		false, // sendRegistryAuth,
	)
}

func (configYml *DockercliLocalConfigConfigWrapperYml) RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions {
	return handler_dockercli_stack_imported.New_RemoveOptions(
		"", // namespace,
	)
}

func (configYml *DockercliLocalConfigConfigWrapperYml) IO() (io.ReadCloser, io.Writer, io.Writer) {
	return nil, nil, nil
}

/**
 * Methods used to load the config yml and conver it to provide settings
 */

func (configYml *DockercliLocalConfigConfigWrapperYml) safe() {
	if &configYml.config == nil {
		if err := configYml.Load(); err != nil {
			log.WithError(err).Error("Could not load dockercli configuration")
		}
	}
}

// Retrieve values by parsing bytes from the wrapper
func (configYml *DockercliLocalConfigConfigWrapperYml) Load() error {
	configYml.config = dockercliLocalConfigureYML{} // reset stored config so that we can repopulate it.

	if sources, err := configYml.configWrapper.Get(CONFIG_KEY_DOCKERCLI_LOCAL); err == nil {
		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedValues := dockercliLocalConfigureYML{} // temporarily hold all config for a specific scope in this
			if err := yaml.Unmarshal(scopedSource, &scopedValues); err == nil {
				configYml.config = scopedValues
				log.WithFields(log.Fields{"bytes": string(scopedSource), "values": scopedValues, "config": configYml}).Debug("Dockercli-Local:Configuration->Load()")
				break
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Couldn't marshall yml scope")
			}
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading dockercli config using key " + CONFIG_KEY_DOCKERCLI_LOCAL)
		return err
	}
}

// Save the current values to the wrapper
func (configYml *DockercliLocalConfigConfigWrapperYml) Save() error {
	/**
	 * @TODO THIS
	 */
	return errors.New("DockercliLocalConfigConfigWrapperYml Set operation not yet written.")
}

/**
 * YML structs
 */

// Wrapper YML struct for all components that could in the yml file
type dockercliLocalConfigureYML struct {
	DeployOptions dockercliLocalConfigureYML_DeployOptions `yml:"Deploy"`
}

// YML holding struct for deploy options, mainly used for the stack handler deploy orchestration
// See github.com/docker/docker/cli/command/stack  (deploy.go) for more understanding
type dockercliLocalConfigureYML_DeployOptions struct {
	Bundlefile       string `yaml:"Bundlefile"`
	Composefile      string `yaml:"Composefile"`
	Namespace        string `yaml:"Namespace"`
	SendRegistryAuth bool   `yaml:"SendRegistryAuth"`
}
