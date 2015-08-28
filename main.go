package main

import (
	"flag"
	"fmt"
	"time"

	"./awsecs"
	"./task"
	log "github.com/Sirupsen/logrus"
)

var conf = flag.String("conf", "", "ECS service family at task definition")
var mode = flag.String("mode", "", "deploy mode at kangol")

func main() {
	finished := make(chan bool)
	go loading(finished)

	flag.Parse()
	deployment, taskDefinition, err := task.ReadConfig(*conf)

	if err != nil {
		log.Fatal(err.Error())
	}

	service := deployment.Service
	cluster := deployment.Cluster
	count := deployment.Count
	oldRevision, err := awsecs.GetOldRevision(service, cluster)

	log.Info("Now Revision is ... ", oldRevision)
	revision := ""

	if *mode == "debug" {
		log.Info("Stop All Tasks at debug mode ....")
		stopTaskError := awsecs.UpdateService(service, cluster, oldRevision, 0)
		if stopTaskError != nil {
			log.Fatal("Stop All Tasks Error -> ", stopTaskError.Error())
		}
		_, err := awsecs.PollingDeployment(service, cluster)
		if err != nil {
			log.Fatal("Stop All Tasks Error -> ", err.Error())
		}
	}

	if *conf != "" {
		newRevision, err := awsecs.RegisterTaskDefinition(taskDefinition)
		if err != nil {
			log.Fatal(err.Error())
		}
		revision = newRevision
	} else {
		revision = oldRevision
	}

	log.Info("Deploying Revision is ... ", revision)
	log.Info("Deploy Start ....")

	getRevisionError := awsecs.UpdateService(service, cluster, revision, count)
	if getRevisionError != nil {
		log.Fatal("UpdateService Error -> ", getRevisionError.Error())
	}

	_, deployError := awsecs.PollingDeployment(service, cluster)
	if deployError != nil {
		rollback := awsecs.UpdateService(service, cluster, oldRevision, count)
		if rollback != nil {
			log.Fatal("Deploy Error & RollBack Revision Error -> ", getRevisionError.Error())
		} else {
			log.Info("RollBack Revision -> ", oldRevision)
			log.Fatal("Deploy Error -> ", deployError.Error())
		}
	} else {
		log.Info("Deploy SUCCESS -> ", service)
	}
	finished <- true

}

func loading(finished chan bool) {
loop:
	for {
		select {
		case <-finished:
			break loop
		default:
			array := []string{"|", "/", "-", "\\"}
			for _, v := range array {
				fmt.Printf("%s\033[1D", v)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}
}
