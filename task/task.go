package task

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"

	"gopkg.in/yaml.v2"
)

type TaskDefinition struct {
	Name                 string
	ContainerDefinitions ContainerDefinition
	Volumes              Volume
}

type ContainerDefinition struct {
	CPU          int64         `yaml:"cpu"`
	Essential    bool          `yaml:"essential"`
	Image        string        `yaml:"image"`
	Memory       int64         `yaml:"memory"`
	PortMappings []PortMapping `yaml:"portMappings"`
	Command      []string      `yaml:"command"`
	EntryPoint   []string      `yaml:"entrypoint"`
	Environment  []Environment `yaml:"environment"`
	Link         []string      `yaml:"link"`
	MountPoints  []MountPoint  `yaml:"mountPoint"`
	VolumesFrom  []VolumesFrom `yaml:"volumesFrom"`
	Volumes      []Volume      `yaml:"volumes"`
}

type PortMapping struct {
	ContainerPort int64  `yaml:"containerPort"`
	HostPort      int64  `yaml:"hostPort"`
	Protocol      string `yaml:"protocol"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type MountPoint struct {
	ContainerPath string `yaml:"containerPath"`
	ReadOnly      bool   `yaml:"readOnly"`
	SourceVolume  string `yaml:"souceVolume"`
}

type VolumesFrom struct {
	ReadOnly        bool   `yaml:"readOnly"`
	SourceContainer string `yaml:"souceContainer"`
}

type Volume struct {
	Host TaskVolumeHost `yaml:"host"`
	Name string         `yaml:"name"`
}

type TaskVolumeHost struct {
	SourcePath string `yaml:"sourcePath"`
}

func ReadConfig(familyName string) (*ecs.RegisterTaskDefinitionInput, error) {
	data, readErr := ioutil.ReadFile("./task-definitions/" + familyName + ".yml")

	if readErr != nil {
		return nil, readErr
	}

	containers := map[string]ContainerDefinition{}
	err := yaml.Unmarshal(data, &containers)

	definitions := []*ecs.ContainerDefinition{}
	volumes := []*ecs.Volume{}
	for name, con := range containers {

		def := &ecs.ContainerDefinition{
			CPU:          aws.Int64(con.CPU),
			Essential:    aws.Bool(con.Essential),
			Image:        aws.String(con.Image),
			Memory:       aws.Int64(con.Memory),
			Name:         aws.String(name),
			PortMappings: getPortMapping(con),
			Command:      getCommands(con),
			EntryPoint:   getEntryPoints(con),
			Environment:  getEnvironments(con),
			Links:        getLinks(con),
			MountPoints:  getMountPoints(con),
			VolumesFrom:  getVolumesFrom(con),
		}
		definitions = append(definitions, def)

		volumes = getVolumes(con)
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: definitions,
		Family:               aws.String(familyName),
		Volumes:              volumes,
	}
	return params, err
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
