# Include variables from the .env file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/docker: run docker compose
.PHONY: run/docker
run/docker:
	docker compose up

## build/docker: build and run docker compose
.PHONY: build/docker
build/docker:
	docker compose up --build

## build/docker/recreate: build, recreate and run docker compose
.PHONY: build/docker/recreate
build/docker/recreate:
	docker compose up --build --remove-orphans --force-recreate

## run/auth: run the auth API
.PHONY: run/auth
run/auth:
	@source ./.env && go run ./auth/main.go

## run/timeline: run the timeline API
.PHONY: run/timeline
run/timeline:
	@source ./.env && go run ./timeline/main.go

## run/tweet: run the tweet API
.PHONY: run/tweet
run/tweet:
	@source ./.env && go run ./tweet/main.go

## run/user: run the user API
.PHONY: run/user
run/user:
	@source ./.env && go run ./user/main.go

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# audit: tidy dependencies and format, vet and test all code
.PHONY: audit
DIRS := auth timeline tweet

audit:
	@for dir in $(DIRS); do \
		echo "Running audit in $$dir..."; \
		(cd $$dir && go fmt ./... && go vet ./... && staticcheck ./... && go test -race -vet=off ./...); \
	done
