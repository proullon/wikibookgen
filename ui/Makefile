IMAGE=proullon/wikibookgen-ui
CONTAINER=wikibookgen-ui

all: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

re: build package run logs ## rebuild binaries, repackage docker image, start locally and run integration tests

wait:
	sleep 5

logs:
	docker logs -f ${CONTAINER}

serve: ## start app locally
	ng serve --configuration=local

build: ## build app
	ng build --configuration=local

package: ## package docker image
	@echo "Packaging production"
	ng build --configuration=production
	docker build -t ${IMAGE}:production .

run:
	docker rm -f ${CONTAINER} || true
	docker run -d --name ${CONTAINER} -p 4061:4200 ${IMAGE}:preprod

push:
	docker push ${IMAGE}:production
	docker push ${IMAGE}:preprod
