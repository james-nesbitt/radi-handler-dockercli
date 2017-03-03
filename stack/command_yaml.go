package stack

import (
	"context"
	"errors"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	docker_cli_compose_loader "github.com/docker/docker/cli/compose/loader"
	docker_cli_compose_types "github.com/docker/docker/cli/compose/types"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"

	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_config "github.com/wunderkraut/radi-api/operation/config"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
	handler_dockercli_stack_container "github.com/wunderkraut/radi-handler-dockercli/stack/container"
)

// A CommandConfigConnector that uses a COnfigWrapper interpreted as yml
type DockercliStackCommand_ConfigureYml struct {
	serviceContext ServiceContext
	configWrapper  api_config.ConfigWrapper

	baseCliOp   *handler_dockercli.DockercliOperationBase
	baseStackOp *DockercliStackOperationBase

	config *DockercliStackCommandConfig
}

// Constructor for DockercliStackCommand_ConfigureYml
func New_DockercliStackCommand_ConfigureYml(configWrapper api_config.ConfigWrapper, serviceContext ServiceContext, baseCliOp *handler_dockercli.DockercliOperationBase, baseStackOp *DockercliStackOperationBase) *DockercliStackCommand_ConfigureYml {
	return &DockercliStackCommand_ConfigureYml{
		configWrapper:  configWrapper,
		serviceContext: serviceContext,
		baseCliOp:      baseCliOp,
		baseStackOp:    baseStackOp,
	}
}

// Convert this to a CommandConfigConnector
func (commandYml *DockercliStackCommand_ConfigureYml) CommandConfigConnector() CommandConfigConnector {
	return CommandConfigConnector(commandYml)
}

/**
 * CommandConfigConnector Interface
 */

// List all commands available
func (commandYml *DockercliStackCommand_ConfigureYml) List(parent string) ([]string, error) {
	commandYml.safe()
	order := commandYml.config.Commands.Order()
	return order, nil
}

// Get a Command interface matching a key
func (commandYml *DockercliStackCommand_ConfigureYml) Get(key string) (api_command.Command, error) {
	commandYml.safe()
	comm, err := commandYml.config.Commands.Get(key)
	if err == nil {
		return comm.Command(), nil
	} else {
		return nil, err
	}
}

/**
 * General functionality
 */

// Safe lazy constructor
func (commandYml *DockercliStackCommand_ConfigureYml) safe() {
	if commandYml.config == nil || commandYml.config.Empty() {
		commandYml.Load()
	}
}

// Retrieve values by parsing bytes from the wrapper
func (commandYml *DockercliStackCommand_ConfigureYml) Load() error {
	commandYml.config = New_DockercliStackCommandConfig() // reset stored config so that we can repopulate it.

	if sources, err := commandYml.configWrapper.Get(CONFIG_KEY_COMMAND); err == nil {
		workingDir := commandYml.serviceContext.WorkingDir()
		envs := commandYml.serviceContext.EnvMap()

		for _, scope := range sources.Order() {
			scopedSource, _ := sources.Get(scope)

			scopedConfig := New_DockercliStackCommandConfig()
			if err := yaml.Unmarshal(scopedSource, &scopedConfig); err == nil {

				/**
				 * This is fun.  Because the docker compose loader has very little exposed
				 * functionality, we actually have to simulate an entire compose file format
				 * even though we want only to have service interpretations.
				 * To do this we sort of jam Commands into services in our yml, and then
				 * pass that to the loader.
				 *
				 * Then when we have service configs, we iterate through the above struct
				 * commands to try to correlate these service configs with the parsed commands.
				 *
				 * It's a bit hackish, but our other options were to write custom parsing, or
				 * to import a massive amount of code (the docker compose loader and all of
				 * it's vendor dependencies.)
				 */

				if scopedLoaderConfig, err := LoadCommandConfig(scopedSource, workingDir, envs); err == nil {
					for _, service := range scopedLoaderConfig.Services {
						id := service.Name
						if comm, err := scopedConfig.Commands.Get(id); err == nil {
							commandYml.serviceContext.AlterService(&service)

							comm.baseCliOp = commandYml.baseCliOp
							comm.baseStackOp = commandYml.baseStackOp
							comm.setServiceConfig(service)

							scopedConfig.Commands.Set(id, comm)
							log.WithFields(log.Fields{"bytes": string(scopedSource), "comm": comm, "id": id}).Debug("Dockercli-Stack:Command->Load(): Added Service to Command")
						}
					}
				} else {
					log.WithError(err).WithFields(log.Fields{"scope": scope}).Warn("Dockercli-Stack:Command->Load(): Failed to add command service config")
				}

			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Warn("Dockercli-Stack:Command->Load(): Failed to unmarshall command source")
			}

			commandYml.config.Merge(*scopedConfig)
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading dockercli config using key " + CONFIG_KEY_COMMAND)
		return err
	}
}

// Retrieve values by parsing bytes from the wrapper
func (commandYml *DockercliStackCommand_ConfigureYml) Save() error {
	return errors.New("Dockercli stack command save not yet written.")
}

// Convert some yaml to a ConfigDetails struct for Loading
func LoadCommandConfig(sourceYmlBytes []byte, workingDir string, envs map[string]string) (*docker_cli_compose_types.Config, error) {
	sourceDict, err := docker_cli_compose_loader.ParseYAML(sourceYmlBytes)
	if err != nil {
		return &docker_cli_compose_types.Config{}, err
	}

	// @todo this is very risky and needs sanity testing X3
	sourceCommands := sourceDict["Commands"].(docker_cli_compose_types.Dict)
	commandMap := map[string]interface{}{}
	for id, command := range sourceCommands {
		commandDict := command.(docker_cli_compose_types.Dict)
		if commandRunDict, exists := commandDict["Run"]; exists {
			commandMap[id] = commandRunDict.(docker_cli_compose_types.Dict)
		}
	}

	sourceDictMap := map[string]interface{}{}
	sourceDictMap["version"] = "3"
	sourceDictMap["services"] = docker_cli_compose_types.Dict(commandMap)

	fakeFile := docker_cli_compose_types.ConfigFile{
		Config: docker_cli_compose_types.Dict(sourceDictMap),
	}
	configDetails := docker_cli_compose_types.ConfigDetails{
		WorkingDir:  workingDir,
		ConfigFiles: []docker_cli_compose_types.ConfigFile{fakeFile},
		Environment: envs,
	}

	if config, err := docker_cli_compose_loader.Load(configDetails); err == nil {
		return config, nil
	} else {
		return &docker_cli_compose_types.Config{}, err
	}
}

/**
 * Structs to holds the YAML results
 */

// YML holding struct for all Commands config
type DockercliStackCommandConfig struct {
	Commands DockercliStackCommands `yaml:"Commands"`
}

// Constructor for DockercliStackCommandConfig
func New_DockercliStackCommandConfig() *DockercliStackCommandConfig {
	return &DockercliStackCommandConfig{}
}

// Is the Config Empty()
func (commConfig *DockercliStackCommandConfig) Empty() bool {
	return &commConfig.Commands == nil || commConfig.Commands.Empty()
}

// Is the Config Empty()
func (commConfig *DockercliStackCommandConfig) Merge(merge DockercliStackCommandConfig) {
	// merge commands
	commConfig.Commands.Merge(merge.Commands)
}

// YML Struct for Commands
type DockercliStackCommands struct {
	comms map[string]CommandYmlCommand
	order []string
}

// Constructor for DockercliStackCommands
func New_DockercliStackCommands() *DockercliStackCommands {
	return &DockercliStackCommands{}
}

// Yaml UnMarshaller
func (comms *DockercliStackCommands) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var holder map[string]CommandYmlCommand
	if error := unmarshal(&holder); error == nil {
		for id, comm := range holder {
			comm.setId(id)
			comms.Set(id, comm)
		}
	}
	return nil
}

// Safe lazy constructor
func (comms *DockercliStackCommands) safe() {
	if comms.comms == nil {
		comms.comms = map[string]CommandYmlCommand{}
		comms.order = []string{}
	}
}

// Safe lazy constructor
func (comms *DockercliStackCommands) Empty() bool {
	return (&comms.comms == nil) || (len(comms.comms) == 0)
}

// Add a command
func (comms *DockercliStackCommands) Set(key string, comm CommandYmlCommand) error {
	comms.safe()
	if _, exists := comms.comms[key]; !exists {
		comms.order = append(comms.order, key)
	}
	comms.comms[key] = comm
	return nil
}

// Get a comm
func (comms *DockercliStackCommands) Get(key string) (CommandYmlCommand, error) {
	comms.safe()
	if com, found := comms.comms[key]; found {
		return com, nil
	} else {
		return com, errors.New("Command not found")
	}
}

// Comm order
func (comms *DockercliStackCommands) Order() []string {
	comms.safe()
	return comms.order
}

// Comm merge
func (comms *DockercliStackCommands) Merge(merge DockercliStackCommands) error {
	comms.safe()
	for _, key := range merge.Order() {
		if _, err := comms.Get(key); err != nil {
			mergeComm, _ := merge.Get(key)
			comms.Set(key, mergeComm)
		}
	}
	return nil
}

type CommandYmlCommand struct {
	scope string
	id    string

	label       string
	description string
	help        string

	disabled   bool
	persistant bool
	internal   bool

	commandProps api_property.Properties

	serviceConfig docker_cli_compose_types.ServiceConfig

	baseCliOp   *handler_dockercli.DockercliOperationBase
	baseStackOp *DockercliStackOperationBase
}

// Yaml UnMarshaller
func (comm *CommandYmlCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var holder struct {
		Scope       string `yaml:"Scope,omitempty"`
		Id          string `yaml:"Id,omitempty"`
		Label       string `yaml:"Label,omitempty"`
		Description string `yaml:"Description,omitempty"`
		Help        string `yaml:"Man,omitempty"`

		Persistant bool `yaml:"Persistant,omitempty"`
		Internal   bool `yaml:"Internal,omitempty"`
		Disabled   bool `yaml:"Disabled,omitempty"`
	}
	if error := unmarshal(&holder); error == nil {
		comm.id = holder.Id
		comm.description = holder.Description
		comm.scope = holder.Scope
		comm.persistant = holder.Persistant
		comm.internal = holder.Internal
		comm.disabled = holder.Disabled
	}

	// @TODO can we get properties from the commands yml?

	return nil
}

// Turn this CommandYmlCommand into a command.Command
func (ymlCommand *CommandYmlCommand) Command() api_command.Command {
	return api_command.Command(ymlCommand)
}

// Assign string Id
func (ymlCommand *CommandYmlCommand) setId(id string) {
	ymlCommand.id = id
}

// Assign string Id
func (ymlCommand *CommandYmlCommand) setServiceConfig(serviceConfig docker_cli_compose_types.ServiceConfig) {
	ymlCommand.serviceConfig = serviceConfig
}

// Return string Scope
func (ymlCommand *CommandYmlCommand) Scope() string {
	return ymlCommand.scope
}

/**
 * Command interace
 */

func (ymlCommand *CommandYmlCommand) Validate() api_result.Result {
	res := api_result.New_StandardResult()

	if ymlCommand.disabled {
		res.AddError(errors.New("Command marked disabled"))
		res.MarkFailed()
	}
	if &ymlCommand.serviceConfig == nil {
		res.AddError(errors.New("No ServiceConfig was defined for the command"))
		res.MarkFailed()
	}

	res.MarkFinished()
	return res.Result()
}

func (ymlCommand *CommandYmlCommand) Usage() api_usage.Usage {
	if ymlCommand.internal {
		return api_operation.Usage_Internal()
	} else {
		return api_operation.Usage_External()
	}
}

// Return string Id
func (ymlCommand *CommandYmlCommand) Id() string {
	return ymlCommand.id
}

// Return string Label
func (ymlCommand *CommandYmlCommand) Label() string {
	return ymlCommand.label
}

// Return string Description
func (ymlCommand *CommandYmlCommand) Description() string {
	return ymlCommand.description
}

// Return string man page
func (ymlCommand *CommandYmlCommand) Help() string {
	return ymlCommand.help
}

// Return string Description
func (ymlCommand *CommandYmlCommand) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	// @TODO find a way to add more dynamic properties from YAML
	if ymlCommand.commandProps != nil {
		props.Merge(ymlCommand.commandProps)
	}

	// make sure that we have a command flags property
	props.Add(api_property.Property(&api_command.CommandFlagsProperty{}))

	return props.Properties()
}

func (ymlCommand *CommandYmlCommand) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	flags := []string{}
	if propFlags, found := props.Get(api_command.OPERATION_PROPERTY_COMMAND_FLAGS); found {
		flags = propFlags.Get().([]string)
	}

	runOpts := handler_dockercli_stack_container.RunOptions{
		name: ymlCommand.Id(),
	}

	// @TODO GET this from a property
	runContext := context.Background()

	dockerCli := ymlCommand.baseCliOp.DockerCli()
	contConfig, hostConfig, netConfig := ConvertServiceToContainerTypes(ymlCommand.serviceConfig)
	log.Printf("SERVICE", ymlCommand.serviceConfig, runOpts, flags, runContext, dockerCli, contConfig, hostConfig, netConfig)

	// RunRun(dockerCli*command.DockerCli, opts*RunOptions, config*container.Config, hostConfig*container.HostConfig, networkingConfig*networktypes.NetworkingConfig)

	// runOptions := libCompose_project_options.Run{
	// 	Detached: false,
	// }

	// // get the service for the command
	// service := ymlCommand.serviceConfig

	// // create a libcompose project
	// project, _ := MakeComposeProject(ymlCommand.projectProperties)

	// // allow our app to alter the service, to do some string replacements etc
	// project.AlterService(&service)

	// project.AddConfig(ymlCommand.Id(), &service)
	// project.Run(runContext, ymlCommand.Id(), flags, runOptions)

	// if !ymlCommand.persistant {
	// 	deleteOptions := libCompose_project_options.Delete{
	// 		RemoveVolume: true,
	// 	}
	// 	project.Delete(runContext, deleteOptions, ymlCommand.Id())
	// }

	res.MarkSuccess()
	res.MarkFinished()

	return res.Result()
}
