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
}

type ContainerDefinition struct {
	CPU          int64         `yaml:"cpu"`
	Essential    bool          `yaml:"essential"`
	Image        string        `yaml:"image"`
	Memory       int64         `yaml:"memory"`
	PortMappings []PortMapping `yaml:"portMappings"`
}

type PortMapping struct {
	ContainerPort int64  `yaml:"containerPort"`
	HostPort      int64  `yaml:"hostPort"`
	Protocol      string `yaml:"protocol"`
}

func ReadConfig(familyName string) (*ecs.RegisterTaskDefinitionInput, error) {
	data, _ := ioutil.ReadFile("./task-definitions/nginx.yml")
	m := map[string]ContainerDefinition{}
	err := yaml.Unmarshal(data, &m)

	definitions := []*ecs.ContainerDefinition{}
	ports := []*ecs.PortMapping{}
	for name, container := range m {
		con := container
		for _, k := range container.PortMappings {
			port := &ecs.PortMapping{
				ContainerPort: aws.Int64(k.ContainerPort),
				HostPort:      aws.Int64(k.HostPort),
				Protocol:      aws.String(k.Protocol),
			}
			ports = append(ports, port)
		}
		def := &ecs.ContainerDefinition{
			CPU:          aws.Int64(con.CPU),
			Essential:    aws.Bool(con.Essential),
			Image:        aws.String(con.Image),
			Memory:       aws.Int64(con.Memory),
			Name:         aws.String(name),
			PortMappings: ports,
		}
		definitions = append(definitions, def)
	}

	params := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: definitions,
		Family:               aws.String(familyName),
	}
	return params, err
}
