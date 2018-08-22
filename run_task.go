package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/recruit-mp/kangol/awscloudwatchlogs"
	"github.com/recruit-mp/kangol/awsecs"
	"github.com/recruit-mp/kangol/task"
)

func runTask(conf, tag string, command string, cpu int64, memory int64) {

	deployment, taskDefinition, err := task.ReadConfig(conf, appendTags(tag))

	if err != nil {
		log.Fatal(err.Error())
	}

	service := deployment.Service
	cluster := deployment.Cluster

	oldRevision, err := awsecs.GetOldRevision(service, cluster)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Now Revision is ... ", oldRevision)
	revision := ""

	if conf != "" {
		newRevision, err := awsecs.RegisterTaskDefinition(taskDefinition)
		if err != nil {
			log.Fatal(err.Error())
		}
		revision = newRevision
	} else {
		revision = oldRevision
	}

	log.Info("Running Revision is ... ", revision)

	commands := []*string{}
	for _, v := range strings.Split(command, " ") {
		commands = append(commands, aws.String(v))
	}

	taskArn, runTaskError := awsecs.RunOneShotTask(cluster, revision, commands, cpu, memory)
	if taskArn != "" {
		log.Info("Show Cloudwatch Logs")
		for _, v := range taskDefinition.ContainerDefinitions {
			logGroup := v.LogConfiguration.Options["awslogs-group"]
			logPrefix := v.LogConfiguration.Options["awslogs-stream-prefix"]
			taskID := strings.Split(taskArn, "/")[1]
			logStream := fmt.Sprintf("%s/%s/%s", *logPrefix, service, taskID)
			events, err := awscloudwatchlogs.GetLogEvents(*logGroup, logStream)
			if err != nil {
				log.Warn("Failed Get Log Events -> ", err.Error())
			}
			for _, v := range events {
				log.Info(v)
			}
		}
	}
	if runTaskError != nil {
		log.Fatal("RunTask Error -> ", runTaskError.Error())
	}
	log.Info("Task Exit Successfully -> ", taskArn)
}
