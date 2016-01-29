# kangol

Container Deploy and Manage System at Amazon EC2 Container Service(ECS).
You can use YAML files to manage docker container(task-definition at ECS).
And when kangol catches the errors at ECS deploy, it can revert before revision automatically.


![](https://cloud.githubusercontent.com/assets/2541396/10562719/ee2773fc-75a4-11e5-8b46-9273ff110db2.gif)

## Installation

```
go get -u github.com/kohey18/kangol
```

## Usage

* setting up environment variable for AWS credential

```
export AWS_ACCESS_KEY_ID=<YOUR AWS_ACCESS_KEY_ID>
export AWS_SECRET_ACCESS_KEY=<YOUR AWS_SECRET_ACCESS_KEY>
export AWS_DEFAULT_REGION=<ECS AWS_DEFAULT_REGION>
```

and then, if you colud through the GOPATH at your environment.
You can use like this.

```
kangol --conf ./nginx.yml
```

## Options


```
GLOBAL OPTIONS:
   --conf               ECS service family at task definition
   --tag                --tag has a container tag
   --debug              --debug has a debug mode
   --loading            --loading has a loading while deploying
   --help, -h           show help
   --version, -v        print the version
```

### tag

if you want to change docker container tag

```
kangol --conf ./api.yml --tag api:<CONTAINER_TAG_NUM>
```

(ex)

```
kangol --conf ./api.yml --tag api:b34cda71
```

### debug mode

You want to stop all container at a service before you want to deployed because ECS show 'resouce not found' because of few your resouces in your cluster.

I'll recommend to apply your development enviroment at ECS.

(ex)

```
kangol --conf ./nginx.yml --tag nginx:latest --debug
```

### loading mode

if you think kangol loading disturbed your deployment enviroment, you should add loading option.
When you off lading option, kangol will not show loading while deploying.

(ex) loading on

```
kangol --conf ./nginx.yml --tag nginx:latest --debug --loading
```

(ex) loading off

```
kangol --conf ./nginx.yml --tag nginx:latest --debug
```

### Setting files

```
cluster: "CLUSTER NAME"
service: "SERVICE NAMAE"
desiredCount: "desirecCount at Cluster"
name: "task-definition-Name"
task:
  "container-name":
    image: "your-registry/image-name:tag"
    cpu: container-cpu-unit
    memory: container-memory
    portMappings:
      - containerPort: container-port
        hostPort: host-port
        protocol: "tcp" etc
    essential: true or false
    environment:
        - name: "JAVA_OPTS"
          value: "-Dconfig.resource=api.conf
    logConfiguration:
        logDriver: fluentd
        options:
            fluentdAddress: "your logger host"
            tag: "docker.{{.Name}}"
```

* example

```
cluster: "dev"
service: "api"
desiredCount: 1
name: "api"
task:
  api:
    image: "kohey18/api:latest"
    cpu: 512
    memory: 1800
    portMappings:
      - containerPort: 9000
        hostPort: 4000
        protocol: "tcp"
    essential: true
    environment:
        - name: "JAVA_OPTS"
          value: "-Dconfig.resource=api.conf
    logConfiguration:
        logDriver: fluentd
        options:
            fluentdAddress: fluentd.example.com
            tag: "docker.{{.Name}}"
  nginx:
    image: "kohey18/nginx:latest"
    cpu: 512
    memory: 1800
    portMappings:
      - containerPort: 9000
        hostPort: 4000
        protocol: "tcp"
    essential: true
```
