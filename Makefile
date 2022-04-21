include version

DOCKER_IMG_TAG := chung1905/gokuwiki:${GOKUWIKI_VERSION}

all: docker_build docker_push

docker_build:
	docker build -t $(DOCKER_IMG_TAG) .

docker_push:
	docker push $(DOCKER_IMG_TAG)
