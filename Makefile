DOCKER_SERVICE_GO ?= go
DOCKER_SERVICE_SAM ?= sam
DOCKER_COMPOSE_RUN = docker-compose run --rm

FUNCTION_NAME ?= OnConnectFunction
STACK_NAME ?= meeting-room-notify-stack
TEST_TARGET ?= ./...

all:
	$(MAKE) install

.env: .env.sample
	cp .env.sample .env
	@echo "An .env file has been created. Please correct."

docker-build: .env
	docker-compose build

install: docker-build
	$(MAKE) build

.PHONY: build
build lint clean:
	$(DOCKER_COMPOSE_RUN) $(DOCKER_SERVICE_GO) make -f Makefile_go $@

test:
	$(DOCKER_COMPOSE_RUN) $(DOCKER_SERVICE_GO) make -f Makefile_go $@ TEST_TARGET=$(TEST_TARGET)

SAM_COMMAND := deploy \
	package \
	validate \
	describe-stack

.PHONY: deploy
$(SAM_COMMAND):
	$(DOCKER_COMPOSE_RUN) $(DOCKER_SERVICE_SAM) make -f Makefile_sam $@

deploy-guided:
	$(DOCKER_COMPOSE_RUN) $(DOCKER_SERVICE_SAM) make -f Makefile_sam $@ \
		GUIDED=--guided

log:
	$(DOCKER_COMPOSE_RUN) $(DOCKER_SERVICE_SAM) make -f Makefile_sam $@ \
		FUNCTION_NAME=$(FUNCTION_NAME)
