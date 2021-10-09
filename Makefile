all: build

build:
	@DOCKER_BUILDKIT=1 docker build . --target bin \
	--output bin/

stop:
	@docker-compose -p boardsite down

development:
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 \
	docker-compose -p boardsite up --abort-on-container-exit --build

.PHONY: all build stop development
