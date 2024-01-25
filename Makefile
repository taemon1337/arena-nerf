IMAGE ?= taemon1337/serf-cluster
VERSION ?= 0.0.1

build:
	go build .

docker-build:
	docker build -t ${IMAGE}:${VERSION} .

docker-push:
	docker push ${IMAGE}:${VERSION}

docker-up:
	docker compose up
