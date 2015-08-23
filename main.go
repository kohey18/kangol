package main

import (
	"flag"
	"fmt"

	"./awsecs"
)

var service = flag.String("service", "", "ECS service name at cluster")
var family = flag.String("family", "", "ECS service family at taskdefinition")
var cluster = flag.String("cluster", "", "ECS cluster name")
var desiredCount = flag.Int64("desiredCount", 1, "desireCount at ECS Service")

func main() {
	flag.Parse()
	result, _ := awsecs.GetOldRevision(*service, *cluster)
	fmt.Printf("[INFO] now revision ---> %s\n", result)
	newRevision, err := awsecs.RegisterTaskDefinition(*family)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("[INFO] GET newRevision ---> %s\n", newRevision)

	getRevisionError := awsecs.UpdateSerive(*service, *cluster, newRevision, *desiredCount)
	if getRevisionError != nil {
		fmt.Printf("[ERROR] UpdateService ERROR ---> %s\n", getRevisionError.Error())
	}

	_, deployError := awsecs.PollingDeployment(*service, *cluster)
	if deployError != nil {
		fmt.Printf("[ERROR] Deploy Error ---> %s\n", deployError.Error())
	} else {
		fmt.Printf("[INFO] service: %s deploy done\n", *service)
	}
}
