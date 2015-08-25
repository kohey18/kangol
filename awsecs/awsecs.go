package awsecs

import (
	"errors"
	"strings"

	"../task"
	log "github.com/Sirupsen/logrus"
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
var deploymentMessage = ""

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
		log.Fatal("RegisterTaskDefinition Error -> ", err.Error())
	}

	resp3, err := svc.RegisterTaskDefinition(params3)
	if err != nil {
		log.Fatal("RegisterTaskDefinition Error -> ", err.Error())
	}
	return strings.Split(*resp3.TaskDefinition.TaskDefinitionARN, "/")[1], err
}

func UpdateService(service, cluster, revision string, desiredCount int64) error {

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

	message, err := checkResouce(deployment.message)
	if err != nil {
		return message, err
	}

	if deploymentMessage != message {
		log.Info(message)
	}

	if (deployment.primary == deployment.desire) && deployment.active == 0 {
		return deployment.message, nil
	} else {
		deploymentMessage = message
		return PollingDeployment(service, cluster)
	}
}

func checkResouce(message string) (string, error) {
	// TODO: コンテナ配置時におけるメッセージにてエラーを検出
	if strings.Contains(message, "resources could not be found") {
		return message, errors.New("resources could not be found")
	} else {
		return message, nil
	}
}
