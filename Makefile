all: build

build:
	@DOCKER_BUILDKIT=1 docker build . --target bin \
	--output bin/

test:
	@docker run --rm --name unit-test-redis \
	-p6379:6379 -d redis:alpine
	@DOCKER_BUILDKIT=1 docker build . --rm --target unit-test \
	--network=host || (docker stop unit-test-redis && exit 1)
	@docker stop unit-test-redis

stop:
	@docker-compose -p boardsite down

development:
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 COMPOSE_TARGET=development \
	docker-compose -p boardsite up --abort-on-container-exit --build

production:
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 COMPOSE_TARGET=production \
	docker-compose -p boardsite up -d --build

.PHONY: all build test stop clean development production
