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
	/**
	 * @TODO check first for client options from yml somehow
	 */

	return configYml.DockercliLocalConfigDefault.ClientOptions()
}

func (configYml *DockercliLocalConfigConfigWrapperYml) DeployOptions() *handler_dockercli_stack_imported.DeployOptions {
	configYml.safe()

	var opts *handler_dockercli_stack_imported.DeployOptions
	if &configYml.config != nil {
		opts = configYml.config.DeployOptions()
	}
	if opts == nil {
		opts = configYml.DockercliLocalConfigDefault.DeployOptions()
	}
	return opts
}

func (configYml *DockercliLocalConfigConfigWrapperYml) RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions {
	configYml.safe()

	var opts *handler_dockercli_stack_imported.RemoveOptions
	if &configYml.config != nil {
		opts = configYml.config.RemoveOptions()
	}
	if opts == nil {
		opts = configYml.DockercliLocalConfigDefault.RemoveOptions()
	}
	return opts
}

func (configYml *DockercliLocalConfigConfigWrapperYml) PsOptions() *handler_dockercli_stack_imported.PsOptions {
	configYml.safe()

	var opts *handler_dockercli_stack_imported.PsOptions
	if &configYml.config != nil {
		opts = configYml.config.PsOptions()
	}
	if opts == nil {
		opts = configYml.DockercliLocalConfigDefault.PsOptions()
	}
	return opts
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
	Client dockercliLocalConfigureYML_Client `yml:"Client"`
	Deploy dockercliLocalConfigureYML_Deploy `yml:"Deploy"`
}

// Convert this to DeployOptions
func (ymlConfig *dockercliLocalConfigureYML) DeployOptions() *handler_dockercli_stack_imported.DeployOptions {
	return ymlConfig.Deploy.DeployOptions()
}

// Convert this to RemoveOptions
func (ymlConfig *dockercliLocalConfigureYML) RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions {
	return ymlConfig.Deploy.RemoveOptions()
}

// Convert this to PSOptions
func (ymlConfig *dockercliLocalConfigureYML) PsOptions() *handler_dockercli_stack_imported.PsOptions {
	return ymlConfig.Deploy.PsOptions()
}

// YML holding struct for client options,
type dockercliLocalConfigureYML_Client struct {
	Debug    bool   `yaml:"Debug"`
	LogLevel string `yaml:"Debug"`

	Hosts     []string `yaml:"Debug"`
	TLS       bool     `yaml:"Debug"`
	TLSVerify bool     `yaml:"Debug"`

	CAFile             string `yaml:"Debug"`
	CertFile           string `yaml:"Debug"`
	KeyFile            string `yaml:"KeyFile"`
	InsecureSkipVerify bool   `yaml:"InsecureSkipVerify"`

	TrustKey  string `yaml:"TrustKey"`
	ConfigDir string `yaml:"ConfigDir"`
	Version   bool   `yaml:"Version"`
}

// Convert this to DeployOptions
func (ymlClient *dockercliLocalConfigureYML_Client) ClientOptions() *docker_cli_flags.ClientOptions {
	common := docker_cli_flags.NewCommonOptions()
	common.Debug = ymlClient.Debug
	common.TLSVerify = ymlClient.TLSVerify
	common.TrustKey = ymlClient.TrustKey

	if ymlClient.LogLevel != "" {
		common.LogLevel = ymlClient.LogLevel
	}
	if len(ymlClient.LogLevel) > 0 {
		common.Hosts = append(common.Hosts, ymlClient.Hosts...)
	}

	return &(docker_cli_flags.ClientOptions{
		Common:    common,
		ConfigDir: ymlClient.ConfigDir,
		Version:   ymlClient.Version,
	})
}

// YML holding struct for deploy options, mainly used for the stack handler deploy orchestration
// See github.com/docker/docker/cli/command/stack  (deploy.go) for more understanding
type dockercliLocalConfigureYML_Deploy struct {
	Bundlefile       string `yaml:"Bundlefile"`
	Composefile      string `yaml:"Composefile"`
	Namespace        string `yaml:"Namespace"`
	SendRegistryAuth bool   `yaml:"SendRegistryAuth"`
}

// Convert this to DeployOptions
func (ymlDeploy *dockercliLocalConfigureYML_Deploy) DeployOptions() *handler_dockercli_stack_imported.DeployOptions {
	return handler_dockercli_stack_imported.New_DeployOptions(
		ymlDeploy.Bundlefile,
		ymlDeploy.Composefile,
		ymlDeploy.Namespace,
		ymlDeploy.SendRegistryAuth,
	)
}

// Convert this to RemoveOptions
func (ymlDeploy *dockercliLocalConfigureYML_Deploy) RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions {
	return handler_dockercli_stack_imported.New_RemoveOptions(
		ymlDeploy.Namespace,
	)
}

// Convert this to PsOptions
func (ymlDeploy *dockercliLocalConfigureYML_Deploy) PsOptions() *handler_dockercli_stack_imported.PsOptions {
	return handler_dockercli_stack_imported.New_PsOptions(
		false,
		ymlDeploy.Namespace,
		false,
		false,
		"",
	)
}
