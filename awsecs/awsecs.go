package awsecs

import (
	"errors"
	"fmt"
	"strings"

	"../task"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Deployments struct {
	active  int64
	primary int64
	desire  int64
	message string
}

var svc = ecs.New(&aws.Config{Region: aws.String("ap-northeast-1")})

func GetOldRevision(service, cluster string) (revision string, err error) {

	params := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(service),
		},
		Cluster: aws.String(cluster),
	}

	resp, err := svc.DescribeServices(params)
	return strings.Split(*resp.Services[0].TaskDefinition, "/")[1], err
}

func RegisterTaskDefinition(familyName string) (revision string, err error) {

	params3, err := task.ReadConfig(familyName)
	if err != nil {
		fmt.Println(err.Error())
	}

	resp3, err := svc.RegisterTaskDefinition(params3)
	return strings.Split(*resp3.TaskDefinition.TaskDefinitionARN, "/")[1], err
}

func UpdateSerive(service, cluster, revision string, desiredCount int64) error {

	params := &ecs.UpdateServiceInput{
		Service:        aws.String(service),
		Cluster:        aws.String(cluster),
		DesiredCount:   aws.Int64(desiredCount),
		TaskDefinition: aws.String(revision),
	}
	_, err := svc.UpdateService(params)
	return err
}

func DescribeDeployedService(service, cluster string) (Deployments, error) {
	param := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(service),
		},
		Cluster: aws.String(cluster),
	}
	res, err := svc.DescribeServices(param)
	deployment := Deployments{}
	deployment.desire = *res.Services[0].DesiredCount

	for _, v := range res.Services[0].Deployments {
		switch *v.Status {
		case "ACTIVE":
			deployment.active = *v.RunningCount
		case "PRIMARY":
			deployment.primary = *v.RunningCount
		}
	}

	deployment.message = *res.Services[0].Events[0].Message
	return deployment, err
}

func PollingDeployment(service, cluster string) (string, error) {
	deployment, err := DescribeDeployedService(service, cluster)

	if err != nil {
		return deployment.message, err
	}

	// TODO: タスクが連続で変わり続ける場合はdeploy失敗
	// (ex) service nginx has started 1 tasks: task 474be549-f9e0-4aee-bf1b-6fbac8e3b445.
	// TODO: pollingCountを元にdeployTimeOutの実装
	// TODO: 同じ文言のときは、クルクル or 点々で増やす

	if (deployment.primary == deployment.desire) && deployment.active == 0 {
		return deployment.message, nil
	} else if resourceCheck(deployment.message) {
		return deployment.message, errors.New("resources could not be found")
	} else {
		fmt.Println(deployment.message)
		return PollingDeployment(service, cluster)
	}
}

func resourceCheck(message string) bool {
	return strings.Contains(message, "resources could not be found")
}
