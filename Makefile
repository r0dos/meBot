BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

export GO111MODULE=on

.PHONY: build
build:
	@echo "-- building binary"
	go build \
		-o ./bin/mebot \
		./cmd/mebot

.PHONY: run
run:
	@echo "-- run meBot"
	nohup ./bin/mebot &

.PHONY: run_nohup
run_nohup:
	@echo "-- run with nohup"
	nohup ./bin/mebot &

.PHONY: build_and_run
build_and_run: build run