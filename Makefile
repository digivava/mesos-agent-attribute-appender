all: build push

build:
	GOARCH=amd64 GOOS=linux go build -o bin/mesos-slave-attribute-appender
	docker build -t digivava/mesos-slave-attribute-appender:001 .

push:
	docker push digivava/mesos-slave-attribute-appender:001
