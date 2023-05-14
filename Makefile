REGISTRY_REPO=fl64
CONTAINER_NAME=echo
CONTAINER_VER:=$(shell git describe --tags)
CONTAINER_VER := $(if $(CONTAINER_VER),$(CONTAINER_VER),$(shell git rev-parse --short HEAD))

HADOLINT_VER:=v2.12.0-alpine
GOLANGLINT_VER:=v1.50.1

CONTAINER_NAME_TAG=$(REGISTRY_REPO)/$(CONTAINER_NAME):$(CONTAINER_VER)
CONTAINER_NAME_LATEST=$(REGISTRY_REPO)/$(CONTAINER_NAME):latest

NAMESPACE:=echo-test

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
	docker run --rm -v "${PWD}":/app:ro -w /app hadolint/hadolint:$(HADOLINT_VER) hadolint /app/Dockerfile
	docker run --rm -v $(PWD):/app:ro -w /app golangci/golangci-lint:$(GOLANGLINT_VER) golangci-lint run -v --timeout=360s

mkcert:
	mkcert test

helm-dev-install:
	helm upgrade --install echo -n $(NAMESPACE) helm/echo --set image.tag=$(CONTAINER_VER) --set podAnnotations.date="\"$(shell date +%s)\""  -f tmp/values.yaml

helm-dev-uninstall:
	helm uninstall -n $(NAMESPACE) echo

werf:
	werf build --repo fl64/echo --add-custom-tag=latest
