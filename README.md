# kangol

Container Deploy and Manage System at Amazon Elastic Container Service(ECS).
You can use YAML files to manage docker container(task-definition at ECS).
And when kangol catches the errors at ECS deploy, it can revert before revision automatically.

![gif](https://cloud.githubusercontent.com/assets/2541396/10562719/ee2773fc-75a4-11e5-8b46-9273ff110db2.gif)

## Installation

```console
go get -u github.com/recruit-mp/kangol
```

## Usage

* setting up environment variable for AWS credential

```console
export AWS_ACCESS_KEY_ID=<YOUR AWS_ACCESS_KEY_ID>
export AWS_SECRET_ACCESS_KEY=<YOUR AWS_SECRET_ACCESS_KEY>
export AWS_DEFAULT_REGION=<ECS AWS_DEFAULT_REGION>
```

and then, if you colud through the GOPATH at your environment.
You can use like this.

```console
kangol --conf ./nginx.yml
```

* run only one ECS task

```console
kangol --conf ./api.yaml --tag api:12345abc --command "bundle exec rake db:migrate"
```

## Options

```console
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

```console
kangol --conf ./api.yml --tag api:<CONTAINER_TAG_NUM>
```

(ex)

```console
kangol --conf ./api.yml --tag api:b34cda71
```

### debug mode

You want to stop all container at a service before you want to deployed because ECS show 'resource not found' because of few your resources in your cluster.

I'll recommend to apply your development environment at ECS.

(ex)

```console
kangol --conf ./nginx.yml --tag nginx:latest --debug
```

### loading mode

if you think kangol loading disturbed your deployment environment, you should add loading option.
When you off lading option, kangol will not show loading while deploying.

(ex) loading on

```console
kangol --conf ./nginx.yml --tag nginx:latest --debug --loading
```

(ex) loading off

```console
kangol --conf ./nginx.yml --tag nginx:latest --debug
```

### Setting files

```yaml
cluster: "CLUSTER NAME"
service: "SERVICE NAMAE"
desiredCount: "desiredCount at Cluster"
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

```yaml
cluster: "dev"
service: "api"
desiredCount: 1
name: "api"
task:
  api:
    image: "recruit-mp/api:latest"
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
    image: "recruit-mp/nginx:latest"
    cpu: 512
    memory: 1800
    portMappings:
      - containerPort: 9000
        hostPort: 4000
        protocol: "tcp"
    essential: true
```

## Usage
- Build
    - make build
- Test
    - make test
- Update dependency 
    - make update
- Update dependency and test
    - make update_test