all: build

build:
	@DOCKER_BUILDKIT=1 docker build . --target bin \
	--output bin/

test:
	@DOCKER_BUILDKIT=1 docker build . --rm --target unit-test

stop:
	@docker-compose -p boardsite down

development:
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 COMPOSE_HOST_PORT=8000 \
	docker-compose -p boardsite up --abort-on-container-exit --build

deploy:
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 COMPOSE_HOST_PORT=80 \
	docker-compose -p boardsite up -d --build

.PHONY: all build test stop development production
