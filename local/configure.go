package local

import (
	"io"
	"os"
	"path"

	docker_cli_flags "github.com/docker/docker/cli/flags"

	api_setting "github.com/wunderkraut/radi-api/operation/setting"

	handler_dockercli_stack_imported "github.com/wunderkraut/radi-handler-dockercli/stack/stack"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

type DockercliLocalConfig interface {
	/**
	 * Elements required for top level github.com/wunderkraut/radi-handler-dockercli
	 */

	// Get Dockercli Client options
	ClientOptions() *docker_cli_flags.ClientOptions
	// Get input and output configurations for the Docker CLI
	IO() (in io.ReadCloser, out io.Writer, err io.Writer)

	/**
	 * github.com/wunderkraut/radi-handler-dockercli/stack.DockercliStackConfig interface
	 */
	DeployOptions() *handler_dockercli_stack_imported.DeployOptions
	RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions
}

/**
 * Docker CLI settings for null testing
 */

type DockercliLocalConfigNull struct{}

// return this struct as a DockercliLocalConfig interface
func (nullsettings *DockercliLocalConfigNull) DockercliLocalConfig() DockercliLocalConfig {
	return DockercliLocalConfig(nullsettings)
}

func (nullsettings *DockercliLocalConfigNull) ClientOptions() *docker_cli_flags.ClientOptions {
	return docker_cli_flags.NewClientOptions()
}

func (nullsettings *DockercliLocalConfigNull) DeployOptions() *handler_dockercli_stack_imported.DeployOptions {
	return handler_dockercli_stack_imported.New_DeployOptions("", "", "", false)
}

func (nullsettings *DockercliLocalConfigNull) RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions {
	return handler_dockercli_stack_imported.New_RemoveOptions("")
}

func (nullsettings *DockercliLocalConfigNull) IO() (io.ReadCloser, io.Writer, io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}

/**
 * Docker CLI settings based on the local project settings
 */

type DockercliLocalConfigDefault struct {
	settings       handler_local.LocalAPISettings
	settingWrapper api_setting.SettingWrapper
}

// Constructor for DockercliLocalConfigDefault
func New_DockercliLocalConfigDefault(settings handler_local.LocalAPISettings, settingWrapper api_setting.SettingWrapper) *DockercliLocalConfigDefault {
	return &DockercliLocalConfigDefault{
		settings:       settings,
		settingWrapper: settingWrapper,
	}
}

// return this struct as a DockercliLocalConfig interface
func (defaultsettings *DockercliLocalConfigDefault) DockercliLocalConfig() DockercliLocalConfig {
	return DockercliLocalConfig(defaultsettings)
}

func (defaultsettings *DockercliLocalConfigDefault) ClientOptions() *docker_cli_flags.ClientOptions {
	return docker_cli_flags.NewClientOptions()
}

func (defaultsettings *DockercliLocalConfigDefault) DeployOptions() *handler_dockercli_stack_imported.DeployOptions {
	// projectName, err := defaultsettings.settingWrapper.Get("Project")
	// if err != nil {
	// 	projectName = "default"
	// }
	projectName := "default"

	return handler_dockercli_stack_imported.New_DeployOptions(
		"", // bundlefile,
		path.Join(defaultsettings.settings.ProjectRootPath, "docker-compose.yml"), // composefile,
		projectName, // namespace,
		false,       // sendRegistryAuth,
	)
}

func (defaultsettings *DockercliLocalConfigDefault) RemoveOptions() *handler_dockercli_stack_imported.RemoveOptions {
	// projectName, err := defaultsettings.settingWrapper.Get("Project")
	// if err != nil {
	// 	projectName = "default"
	// }
	projectName := "default"

	return handler_dockercli_stack_imported.New_RemoveOptions(
		projectName, // namespace,
	)
}

func (defaultsettings *DockercliLocalConfigDefault) IO() (io.ReadCloser, io.Writer, io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}
