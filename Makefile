all: build

stop:
	@docker-compose -p boardsite down

start:
	@docker-compose -p boardsite up --abort-on-container-exit --build

.PHONY: all stop start
