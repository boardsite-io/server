all: build

build:
	@DOCKER_BUILDKIT=1 docker build . --target bin \
	--output bin/

test:
	@docker run --rm --name unit-test-redis \
	-p6379:63000 -d redis:alpine
	@DOCKER_BUILDKIT=1 docker build . --rm --target unit-test \
	--network=host || docker stop unit-test-redis
	@docker rm -f unit-test-redis

debug:
	@docker run --rm --name b-redis-debug \
	-p6379:6379 -d redis:alpine	
	@DOCKER_BUILDKIT=1 docker build . \
	--target debug -t b-debug
	@docker run --rm --name b-debug-local -p8000:8000 \
	--link b-redis-debug b-debug && docker stop b-redis-debug

.PHONY: all build test