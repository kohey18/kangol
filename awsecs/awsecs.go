package awsecs

import (
	"errors"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var svc = ecs.New(&aws.Config{Region: aws.String("ap-northeast-1")})
var deploymentMessage = ""
var pollingCount = 0

// Deployments has deployment infomation at ECS
type Deployments struct {
	active  int64
	primary int64
	desire  int64
	message string
}

// GetOldRevision can get revision you specified
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

// RegisterTaskDefinition can get register task-definition using your yml file
func RegisterTaskDefinition(taskDefinition *ecs.RegisterTaskDefinitionInput) (revision string, err error) {

	resp, err := svc.RegisterTaskDefinition(taskDefinition)
	if err != nil {
		log.Fatal("RegisterTaskDefinition Error -> ", err.Error())
	}
	return strings.Split(*resp.TaskDefinition.TaskDefinitionArn, "/")[1], err
}

// UpdateService can update service using revision you specified
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

// DescribeDeployedService can get running count
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

// PollingDeployment can check deployment message at update-service
func PollingDeployment(service, cluster string) (string, error) {
	time.Sleep(100 * time.Millisecond)
	deployment, err := DescribeDeployedService(service, cluster)

	if err != nil {
		return deployment.message, err
	}

	// TODO: タスクが連続で変わり続ける場合はdeploy失敗
	// (ex) service nginx has started 1 tasks: task 474be549-f9e0-4aee-bf1b-6fbac8e3b445.
	// TODO: pollingCountを元にdeployTimeOutの実装

	if deploymentMessage == "" {
		deploymentMessage = deployment.message
	} else if deploymentMessage != deployment.message {
		pollingCount = 0
		log.Info(deployment.message)
		_, err := checkResouce(deployment.message)
		if err != nil {
			return deployment.message, err
		}
	} else {
		if pollingCount > 500 {
			return deployment.message, errors.New(deployment.message)
		}
	}

	if (deployment.primary == deployment.desire) && deployment.active == 0 {
		return deployment.message, nil
	}

	deploymentMessage = deployment.message
	pollingCount++
	return PollingDeployment(service, cluster)

}

func checkResouce(message string) (string, error) {
	// TODO: コンテナ配置時におけるメッセージにてエラーを検出
	if strings.Contains(message, "resources could not be found") {
		return message, errors.New("resources could not be found")
	}
	return message, nil
}
