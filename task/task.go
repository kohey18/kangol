package task

import (
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"

	"gopkg.in/yaml.v2"
)

type ClusterService struct {
	Cluster string
	Service string
	Count   int64
}

// Deployment has deployment config setting
type Deployment struct {
	Cluster string                         `yaml:"cluster"`
	Service string                         `yaml:"service"`
	Count   int64                          `yaml:"desiredCount"`
	Name    string                         `yaml:"name"`
	Task    map[string]ContainerDefinition `yaml:"task"`
}

// ContainerDefinition is struct for ECS TaskDefinition
type ContainerDefinition struct {
	CPU              int64            `yaml:"cpu"`
	Essential        bool             `yaml:"essential"`
	Image            string           `yaml:"image"`
	Memory           int64            `yaml:"memory"`
	PortMappings     []PortMapping    `yaml:"portMappings"`
	Command          []string         `yaml:"command"`
	EntryPoint       []string         `yaml:"entrypoint"`
	Environment      []Environment    `yaml:"environment"`
	Link             []string         `yaml:"links"`
	MountPoints      []MountPoint     `yaml:"mountPoint"`
	VolumesFrom      []VolumesFrom    `yaml:"volumesFrom"`
	Volumes          []Volume         `yaml:"volumes"`
	LogConfiguration LogConfiguration `yaml:"logConfiguration"`
}

// PortMapping is struct for ECS TaskDefinition's PortMapping
type PortMapping struct {
	ContainerPort int64  `yaml:"containerPort"`
	HostPort      int64  `yaml:"hostPort"`
	Protocol      string `yaml:"protocol"`
}

// Environment is struct for TaskDefinition's Environment
type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// MountPoint is struct for TaskDefinition's MoutPoint
type MountPoint struct {
	ContainerPath string `yaml:"containerPath"`
	ReadOnly      bool   `yaml:"readOnly"`
	SourceVolume  string `yaml:"souceVolume"`
}

// VolumesFrom is struct for TaskDefinition's VolumesFrom
type VolumesFrom struct {
	ReadOnly        bool   `yaml:"readOnly"`
	SourceContainer string `yaml:"souceContainer"`
}

// Volume is struct for TaskDefinition's Volume
type Volume struct {
	Host TaskVolumeHost `yaml:"host"`
	Name string         `yaml:"name"`
}

// TaskVolumeHost is struct for TaskDefinition's Volumes's volumeHost
type TaskVolumeHost struct {
	SourcePath string `yaml:"sourcePath"`
}

// LogConfiguration is struct for TaskDefinition's LogConfiguration
type LogConfiguration struct {
	LogDriver string                  `yaml:"logDriver"`
	Options   LogConfigurationOptions `yaml:"options"`
}

// LogConfigurationOptions is struct for LogConfiguration's LogConfigurationOptions
type LogConfigurationOptions struct {
	FluentdAddress string `yaml:"fluentdAddress"`
	Tag            string `yaml:"tag"`
}

// ReadConfig can read config yml
func ReadConfig(conf string, tags map[string]string) (ClusterService, *ecs.RegisterTaskDefinitionInput, error) {
	data, readErr := ioutil.ReadFile(conf)

	if readErr != nil {
		return ClusterService{}, nil, readErr
	}

	containers := Deployment{}
	err := yaml.Unmarshal(data, &containers)

	clusterService := ClusterService{containers.Cluster, containers.Service, containers.Count}

	definitions := []*ecs.ContainerDefinition{}
	volumes := []*ecs.Volume{}

	for name, con := range containers.Task {
		if tags[name] != "" {
			con.Image = addImageTag(con.Image, tags[name])
		}

		def := &ecs.ContainerDefinition{
			Cpu:              aws.Int64(con.CPU),
			Essential:        aws.Bool(con.Essential),
			Image:            aws.String(con.Image),
			Memory:           aws.Int64(con.Memory),
			Name:             aws.String(name),
			PortMappings:     getPortMapping(con),
			Command:          getCommands(con),
			EntryPoint:       getEntryPoints(con),
			Environment:      getEnvironments(con),
			Links:            getLinks(con),
			MountPoints:      getMountPoints(con),
			VolumesFrom:      getVolumesFrom(con),
			LogConfiguration: getLogConfiguration(con),
		}
		definitions = append(definitions, def)

		volumes = getVolumes(con)
	}

	taskDefinitions := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: definitions,
		Family:               aws.String(containers.Name),
		Volumes:              volumes,
	}

	return clusterService, taskDefinitions, err
}

func addImageTag(image, tag string) string {
	imagesArray := strings.Split(image, ":")
	imagesArray[len(imagesArray)-1] = tag
	return strings.Join(imagesArray, ":")
}

func getPortMapping(con ContainerDefinition) []*ecs.PortMapping {
	ports := []*ecs.PortMapping{}
	for _, v := range con.PortMappings {
		port := &ecs.PortMapping{
			ContainerPort: aws.Int64(v.ContainerPort),
			HostPort:      aws.Int64(v.HostPort),
			Protocol:      aws.String(v.Protocol),
		}
		ports = append(ports, port)
	}
	return ports
}

func getCommands(con ContainerDefinition) []*string {
	commands := []*string{}
	for _, v := range con.Command {
		commands = append(commands, aws.String(v))
	}
	return commands
}

func getEntryPoints(con ContainerDefinition) []*string {
	entrypoints := []*string{}
	for _, v := range con.EntryPoint {
		entrypoints = append(entrypoints, aws.String(v))
	}
	return entrypoints
}

func getEnvironments(con ContainerDefinition) []*ecs.KeyValuePair {
	environments := []*ecs.KeyValuePair{}
	for _, v := range con.Environment {
		env := &ecs.KeyValuePair{
			Name:  aws.String(v.Name),
			Value: aws.String(v.Value),
		}
		environments = append(environments, env)
	}
	return environments
}

func getLinks(con ContainerDefinition) []*string {
	links := []*string{}
	for _, v := range con.Link {
		links = append(links, aws.String(v))
	}
	return links
}

func getMountPoints(con ContainerDefinition) []*ecs.MountPoint {
	mountPoints := []*ecs.MountPoint{}
	for _, v := range con.MountPoints {
		mountPoint := &ecs.MountPoint{
			ContainerPath: aws.String(v.ContainerPath),
			ReadOnly:      aws.Bool(v.ReadOnly),
			SourceVolume:  aws.String(v.SourceVolume),
		}
		mountPoints = append(mountPoints, mountPoint)
	}
	return mountPoints
}

func getVolumesFrom(con ContainerDefinition) []*ecs.VolumeFrom {
	volumesFroms := []*ecs.VolumeFrom{}
	for _, v := range con.VolumesFrom {
		volumesFrom := &ecs.VolumeFrom{
			ReadOnly:        aws.Bool(v.ReadOnly),
			SourceContainer: aws.String(v.SourceContainer),
		}
		volumesFroms = append(volumesFroms, volumesFrom)
	}
	return volumesFroms
}

func getVolumes(con ContainerDefinition) []*ecs.Volume {
	volumes := []*ecs.Volume{}
	for _, v := range con.Volumes {
		vol := &ecs.Volume{
			Host: &ecs.HostVolumeProperties{
				SourcePath: aws.String(v.Host.SourcePath),
			},
			Name: aws.String(v.Name),
		}
		volumes = append(volumes, vol)
	}
	return volumes
}

func getLogConfiguration(con ContainerDefinition) *ecs.LogConfiguration {
	logConf := con.LogConfiguration
	conf := &ecs.LogConfiguration{
		LogDriver: aws.String(con.LogConfiguration.LogDriver),
		Options: map[string]*string{
			"fluentd-address": aws.String(logConf.Options.FluentdAddress),
			"tag":             aws.String(logConf.Options.Tag),
		},
	}
	return conf
}
