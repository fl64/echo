REGISTRY_REPO=fl64
CONTAINER_NAME=echo
CONTAINER_VER:=$(shell git describe --tags)
CONTAINER_VER := $(if $(CONTAINER_VER),$(CONTAINER_VER),$(shell git rev-parse --short HEAD))

HADOLINT_VER:=v1.22.1
GOLANGLINT_VER:=v1.39.0

CONTAINER_NAME_TAG=$(REGISTRY_REPO)/$(CONTAINER_NAME):$(CONTAINER_VER)
CONTAINER_NAME_LATEST=$(REGISTRY_REPO)/$(CONTAINER_NAME):latest

.PHONY: up down build latest push push_latest lint
up:
	docker-compose up -d --build

down:
	docker-compose down

build:
	docker build -t $(CONTAINER_NAME_TAG) .

latest: build
	docker tag $(CONTAINER_NAME_TAG) $(CONTAINER_NAME_LATEST)

push: build
	docker push $(CONTAINER_NAME_TAG)

push_latest: push latest
	docker push $(CONTAINER_NAME_LATEST)

lint:
	docker run --rm -v $(PWD):/app:ro -w /app golangci/golangci-lint:$(GOLANGLINT_VER) golangci-lint run -v
	docker run --rm -v "${PWD}":/app:ro -w /app hadolint/hadolint:$(HADOLINT_VER) hadolint /app/Dockerfile

mkcert:
	mkcert test