registry=lanbing322
time:=$(shell date +"%Y%m%d-%H%M%S")

.PHONY all: clean build image install clean

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deepblue-k8s-scheduler .

image:
	docker build -t $(registry)/public:$(time) .
	docker push $(registry)/public:$(time)
	docker rmi $(registry)/public:$(time)