BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

export GO111MODULE=on

.PHONY: build
build:
	@echo "-- building binary with sqlite"
	go build \
		-tags=sqlite \
		-o ./bin/mebot \
		./cmd/mebot

.PHONY: build_postgres
build_postgres:
	@echo "-- building binary with postgres"
	go build \
		-tags=postgres \
		-o ./bin/mebot \
		./cmd/mebot

.PHONY: run
run:
	@echo "-- run meBot"
	./bin/mebot

.PHONY: run_nohup
run_nohup:
	@echo "-- run with nohup"
	nohup ./bin/mebot &

.PHONY: build_and_run
build_and_run: build run

.PHONY: docker
docker:
	@echo "-- building docker container"
	docker build -f Dockerfile -t mebot .

.PHONY: docker_run
docker_run:
	@echo "-- starting docker container"
	docker run --name mebot --rm \
	-v $(pwd)/data:/persis \
	--env DB_URL=/persis/me_bot.sqlite \
	-d mebot
#	docker run --name mebot --rm \
# 	-v /Users/r0dos/go/src/meBot/data:/persis \
# 	--env DB_URL=/persis/me_bot.sqlite \
# 	-d mebot
