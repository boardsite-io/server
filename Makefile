all: build

.PHONY: build
build:
	@DOCKER_BUILDKIT=1 docker build . --target bin \
	--output bin/

.PHONY: test
test:
	@docker run --rm --name unit-test-redis \
	-p6379:63790 -d redis:alpine
	@DOCKER_BUILDKIT=1 docker build . --rm --target unit-test \
	--network=host || docker stop unit-test-redis
	@docker rm -f unit-test-redis