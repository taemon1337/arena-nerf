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

controller:
	docker run --rm -it \
    -v ./logs:/tmp/logs:rw \
		--net host \
		${IMAGE}:${VERSION} \
		-name control \
		-role ctrl \
		-server \
		-mode domination \
		-start \
		-allow-api-actions \
		-logdir /tmp/logs \
		-gametime 1m \
		-expect 4 \
		-tag role=ctrl \
		-team blue \
		-team red \
		-team green \
		-team yellow

