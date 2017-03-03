package stack

import (
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/cli/compose/types"

	"github.com/docker/go-connections/nat"
)

// Convert the compose Service to the related container info structs.
func ConvertServiceToContainerTypes(service types.ServiceConfig) (container.Config, container.HostConfig, network.NetworkingConfig) {

	// service :] ServiceConfig {
	// 	Name string

	// 	CapAdd          []string `mapstructure:"cap_add"`
	// 	CapDrop         []string `mapstructure:"cap_drop"`
	// 	CgroupParent    string   `mapstructure:"cgroup_parent"`
	// 	Command         ShellCommand
	// 	ContainerName   string   `mapstructure:"container_name"`
	// 	DependsOn       []string `mapstructure:"depends_on"`
	// 	Deploy          DeployConfig
	// 	Devices         []string
	// 	DNS             StringList
	// 	DNSSearch       StringList `mapstructure:"dns_search"`
	// 	DomainName      string     `mapstructure:"domainname"`
	// 	Entrypoint      ShellCommand
	// 	Environment     MappingWithEquals
	// 	EnvFile         StringList `mapstructure:"env_file"`
	// 	Expose          StringOrNumberList
	// 	ExternalLinks   []string         `mapstructure:"external_links"`
	// 	ExtraHosts      MappingWithColon `mapstructure:"extra_hosts"`
	// 	Hostname        string
	// 	HealthCheck     *HealthCheckConfig
	// 	Image           string
	// 	Ipc             string
	// 	Labels          MappingWithEquals
	// 	Links           []string
	// 	Logging         *LoggingConfig
	// 	MacAddress      string `mapstructure:"mac_address"`
	// 	NetworkMode     string `mapstructure:"network_mode"`
	// 	Networks        map[string]*ServiceNetworkConfig
	// 	Pid             string
	// 	Ports           []ServicePortConfig
	// 	Privileged      bool
	// 	ReadOnly        bool `mapstructure:"read_only"`
	// 	Restart         string
	// 	Secrets         []ServiceSecretConfig
	// 	SecurityOpt     []string       `mapstructure:"security_opt"`
	// 	StdinOpen       bool           `mapstructure:"stdin_open"`
	// 	StopGracePeriod *time.Duration `mapstructure:"stop_grace_period"`
	// 	StopSignal      string         `mapstructure:"stop_signal"`
	// 	Tmpfs           StringList
	// 	Tty             bool `mapstructure:"tty"`
	// 	Ulimits         map[string]*UlimitsConfig
	// 	User            string
	// 	Volumes         []string
	// 	WorkingDir      string `mapstructure:"working_dir"`
	// }

	healthCheck := convertHealthCheck(*service.HealthCheck)

	ContConfig := container.Config{
		Hostname: service.Hostname,
		// Domainname: service.Domainname,
		User: service.User,
		// AttachStdin     bool                // Attach the standard input, makes possible user interaction
		// AttachStdout    bool                // Attach the standard output
		// AttachStderr    bool                // Attach the standard error
		// ExposedPorts:    nat.ParsePortSpecs([]string{}) // @TODO get this from the []ServicePortConfig
		Tty:         service.Tty,
		OpenStdin:   true, // @TODO should I be doing this?
		StdinOnce:   true, // @TODO  ||
		Env:         convertMappingWithEquals(service.Environment),
		Cmd:         convertCommand(service.Command),
		Healthcheck: &healthCheck,
		ArgsEscaped: false,
		Image:       service.Image,
		// Volumes:     service.Volumes, // @TODO will this work? :: map[string]struct{} // List of volumes (mounts) used for the container
		WorkingDir: service.WorkingDir,
		Entrypoint: convertCommand(service.Entrypoint),
		// NetworkDisabled bool                `json:",omitempty"` // Is network disabled
		MacAddress: service.MacAddress,
		// OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
		Labels:     map[string]string(service.Labels),
		StopSignal: service.StopSignal,
		// StopTimeout: convertDuration(service.StopGracePeriod),
		// Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
	}

	HostConfig := container.HostConfig{
		// Applicable to all platforms
		Binds: service.Volumes,
		// ContainerIDFile string        // File (path) where the containerId is written
		// LogConfig       LogConfig     // Configuration of the logs for this container
		// NetworkMode     NetworkMode   // Network mode to use for the container
		// PortBindings: convertPorts(service.Ports),
		// RestartPolicy   RestartPolicy // Restart policy to be used for the container
		// AutoRemove      bool          // Automatically remove container when it exits
		// VolumeDriver    string        // Name of the volume driver used to mount volumes
		// VolumesFrom     []string      // List of volumes to take from other container

		// Applicable to UNIX platforms
		CapAdd:  strslice.StrSlice(service.CapAdd),
		CapDrop: strslice.StrSlice(service.CapDrop),
		DNS:     convertStringList(service.DNS),
		// DNSOptions      []string          `json:"DnsOptions"` // List of DNSOption to look for
		DNSSearch:  convertStringList(service.DNSSearch),
		ExtraHosts: convertMappingWithColon(service.ExtraHosts),
		// GroupAdd        []string          // List of additional groups that the container process will run as
		IpcMode: container.IpcMode(service.Ipc),
		Cgroup:  container.CgroupSpec(service.CgroupParent),
		Links:   service.Links,
		// OomScoreAdj     int               // Container preference for OOM-killing
		PidMode:    container.PidMode(service.Pid),
		Privileged: service.Privileged,
		// PublishAllPorts bool              // Should docker publish all exposed port for the container
		ReadonlyRootfs: service.ReadOnly,
		SecurityOpt:    service.SecurityOpt,
		// StorageOpt      map[string]string `json:",omitempty"` // Storage driver options per container.
		// Tmpfs           map[string]string `json:",omitempty"` // List of tmpfs (mounts) used for the container
		// UTSMode         UTSMode           // UTS namespace to use for the container
		// UsernsMode      UsernsMode        // The user namespace to use for the container
		// ShmSize         int64             // Total shm memory usage
		// Sysctls         map[string]string `json:",omitempty"` // List of Namespaced sysctls used for the container
		// Runtime         string            `json:",omitempty"` // Runtime to use with this container

		// Applicable to Windows
		// ConsoleSize [2]uint   // Initial console size (height,width)
		// Isolation   Isolation // Isolation technology of the container (e.g. default, hyperv)

		// Contains container's resources (cgroups, ulimits)
		// Resources // @TODO THIS?

		// Mounts specs used by the container
		Mounts: convertVolumesToMounts(service.Volumes),

		// Run a custom init inside the container, if null, use the daemon's configured settings
		// Init *bool `json:",omitempty"`

		// Custom init path
		// InitPath string `json:",omitempty"`
	}

	NetworkingEndpoints := network.NetworkingConfig{
	// EndpointsConfig: map[string]*EndpointSettings // Endpoint configs for each connecting network
	}

	return ContConfig, HostConfig, NetworkingEndpoints
}

func convertMappingWithEquals(mapping types.MappingWithEquals) []string {
	envs := []string{}
	for key, value := range map[string]string(mapping) {
		envs = append(envs, strings.Join([]string{key, value}, "="))
	}
	return envs
}

func convertMappingWithColon(mapping types.MappingWithColon) []string {
	envs := []string{}
	for key, value := range map[string]string(mapping) {
		envs = append(envs, strings.Join([]string{key, value}, ":"))
	}
	return envs
}

func convertCommand(com types.ShellCommand) strslice.StrSlice {
	return strslice.StrSlice([]string(com))
}

func convertStringList(list types.StringList) strslice.StrSlice {
	return strslice.StrSlice([]string(list))
}

func convertHealthCheck(test types.HealthCheckConfig) container.HealthConfig {
	return container.HealthConfig{Test: []string{"NONE"}} // @TODO properly handle this
}

func converDuration(dur time.Duration) int {
	return int(dur)
}

func convertVolumesToMounts(volumes []string) []mount.Mount {
	return []mount.Mount{}
}

func convertPorts(ports []types.ServicePortConfig) nat.PortMap {
	return nat.PortMap{}
}
