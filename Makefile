CRDB_HOST=`cat .dev.conf | grep CRDB_HOST | cut -d '=' -f 2`
DOCKER_NAME=`cat .dev.conf | grep DOCKER_NAME | cut -d '=' -f 2`
DOCKER_IMAGE=`cat .dev.conf | grep DOCKER_IMAGE | cut -d '=' -f 2`
FOLDER=`cat .dev.conf | grep FOLDER | cut -d '=' -f 2`
LANGUAGE=`cat .dev.conf | grep LANGUAGE | cut -d '=' -f 2`
DB_NAME=`cat .dev.conf | grep DB_NAME | cut -d '=' -f 2`
MOUNT_POINT=`cat .dev.conf | grep MOUNT_POINT | cut -d '=' -f 2`
EXPOSED_PORT=`cat .dev.conf | grep EXPOSED_PORT | cut -d '=' -f 2`

PROTOC_OPTS ?=		-I/usr/local/include -I./$(PROTO_DIR) -I$(GOPATH) -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
PROTO_TMPL=../genkit
IMPORT_PATH=github.com/proullon/wikibookgen/api
DIR =./api
PROTO_FILE =cmd/wikibookgen/wikibookgen.proto


all: help

generate: ## generate Go from proto file
	rm -f $(shell find . -name "*.gen.go" -not -path "./vendor/*")
	TARGET_PATH="$(IMPORT_PATH)" protoc $(PROTOC_OPTS) --gotemplate_out=debug=true,template_dir=$(PROTO_TMPL):$(DIR) $(PROTO_FILE)
	goimports -w $(shell find . -name "*.gen.go" -not -path "./vendor/*")
	mv api/model/model.gen.ts ui/src/app/wikibook.ts

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

re: generate install test package run wait it logs ## rebuild binaries

install: ## install binaries
	clear
	go build ./cmd/wikibookgen
	go build ./cmd/integration

test:
	go test -short -race ./...

stop: ## stop docker container
	docker rm -f $(DOCKER_NAME) || true

run: stop ## start docker container
	docker run --restart=always -d --mount type=bind,source=$(MOUNT_POINT),target=/tmp/wikibookgen --name $(DOCKER_NAME) -p $(EXPOSED_PORT):8080 -e CRDB_HOST=$(CRDB_HOST) $(DOCKER_IMAGE):latest

logs: ## show docker logs
	docker logs -f $(DOCKER_NAME)

package: ## build docker image
	docker build -t $(DOCKER_IMAGE):latest .

it: ## run integration tests
	./integration

wait:
	sleep 10

