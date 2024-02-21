IMAGE ?= taemon1337/arena-nerf
VERSION ?= 0.0.1

build:
	go build .

pibuild:
	GOOS=linux GOARCH=arm64 go build -o arena-nerf.arm64

docker-build:
	docker build -t ${IMAGE}:${VERSION} .

docker-push:
	docker push ${IMAGE}:${VERSION}

docker-up:
	docker compose up
