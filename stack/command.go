package stack

import (
	"errors"

	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"

	handler_dockercli "github.com/wunderkraut/radi-handler-dockercli"
)

/**
 * Implement command containers that mix into
 * Dockercli Stack orchestrated containers
 */

const (
	CONFIG_KEY_COMMAND = "commands" // The Config key for settings
)

// A wrapper interface which pulls command information from a config wrapper backend
type CommandConfigConnector interface {
	List(parent string) ([]string, error)
	Get(key string) (api_command.Command, error)
}

/**
 * Operations
 */

// LibCompose Command List operation
type DockercliStackCommandListOperation struct {
	api_command.BaseCommandListOperation
	api_command.BaseCommandKeyKeysOperation
	handler_dockercli.DockercliOperationBase
	DockercliStackOperationBase

	Connector CommandConfigConnector
}

// Validate the operation
func (list *DockercliStackCommandListOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Get properties
func (list *DockercliStackCommandListOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(list.BaseCommandKeyKeysOperation.Properties())

	return props.Properties()
}

// Execute the Dockercli Stack Command List operation
func (list *DockercliStackCommandListOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	keysProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEYS)

	parent := ""
	if key, ok := keyProp.Get().(string); ok && key != "" {
		parent = key
	}

	if keyList, err := list.Connector.List(parent); err == nil {
		keysProp.Set(keyList)
		res.MarkSuccess()
	} else {
		res.MarkFailed()
		res.AddError(err)
	}

	res.MarkFinished()

	return res.Result()
}

// Dockercli Stack Command Get operation
type DockercliStackCommandGetOperation struct {
	api_command.BaseCommandGetOperation
	api_command.BaseCommandKeyCommandOperation
	handler_dockercli.DockercliOperationBase
	DockercliStackOperationBase

	Connector CommandConfigConnector
}

// Validate the operation
func (get *DockercliStackCommandGetOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Get properties
func (get *DockercliStackCommandGetOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(get.BaseCommandKeyCommandOperation.Properties())

	return props.Properties()
}

// Execute the Dockercli Stack Command Get operation
func (get *DockercliStackCommandGetOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	commandProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_COMMAND)

	if key, ok := keyProp.Get().(string); ok && key != "" {

		if comm, err := get.Connector.Get(key); err == nil {
			// pass all props to make a project
			commandProp.Set(comm)
			res.MarkSuccess()
		} else {
			res.AddError(err)
			res.MarkFailed()
		}

	} else {
		res.AddError(errors.New("No command name provided."))
		res.MarkFailed()
	}

	res.MarkFinished()

	return res.Result()
}
