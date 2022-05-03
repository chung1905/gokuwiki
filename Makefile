include version

DOCKER_IMG_NAME := chung1905/gokuwiki
DOCKER_IMG_TAG := $(DOCKER_IMG_NAME):${GOKUWIKI_VERSION}

all: clean build docker_build docker_push

clean:
	rm -f gokuwiki

build:
	go build .

docker_build:
	docker build -t $(DOCKER_IMG_TAG) -t $(DOCKER_IMG_NAME):latest .

docker_push:
	docker push $(DOCKER_IMG_TAG)
	docker push $(DOCKER_IMG_NAME):latest
