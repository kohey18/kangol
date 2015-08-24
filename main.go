package main

import (
	"flag"
	"fmt"
	"time"

	"./awsecs"
	log "github.com/Sirupsen/logrus"
)

var service = flag.String("service", "", "ECS service name at cluster")
var family = flag.String("family", "", "ECS service family at task definition")
var cluster = flag.String("cluster", "", "ECS cluster name")
var desiredCount = flag.Int64("desiredCount", 1, "desireCount at ECS Service")

func main() {
	finished := make(chan bool)
	go loading(finished)

	flag.Parse()

	result, _ := awsecs.GetOldRevision(*service, *cluster)
	log.Info("Now Revision is ... ", result)

	newRevision, err := awsecs.RegisterTaskDefinition(*family)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("New Revision is ... ", newRevision)
	log.Info("Deploy Start ....")

	getRevisionError := awsecs.UpdateSerive(*service, *cluster, newRevision, *desiredCount)
	if getRevisionError != nil {
		log.Fatal("UpdateService Error -> ", getRevisionError.Error())
	}

	_, deployError := awsecs.PollingDeployment(*service, *cluster)
	if deployError != nil {
		log.Fatal("Deploy Error -> ", deployError.Error())
	} else {
		log.Info("Deploy SUCCESS -> ", *service)
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
