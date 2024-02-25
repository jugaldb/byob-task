#!/usr/bin/env just --justfile

GO := "go"

GOVET_COMMAND := GO + " vet"
GOTEST_COMMAND := GO + " test"
GOCOVER_COMMAND := GO + " tool cover"
GOBUILD_COMMAND := GO + " build"
COVTHRESHOLD := "95"


# display all commands
default:
    @just --list --unsorted

# Run static checks
check:
    {{GOVET_COMMAND}} ./...

# Execute test cases with code coverage
test:
    {{GOTEST_COMMAND}} -v -race -covermode=atomic -coverprofile=coverage.out ./...
    @{{GOCOVER_COMMAND}} -func=coverage.out
    @{{GOCOVER_COMMAND}} -html=coverage.out -o coverage.html

# Clean dist directory and rebuild the binary file
build:
    rm -rf ./dist && CGO_ENABLED=0 {{GOBUILD_COMMAND}} -ldflags="-w -s" -o ./dist/app ./src

debug:
    rm -rf ./debugApp && CGO_ENABLED=0 {{GOBUILD_COMMAND}} -gcflags="all=-N -l" -o debugApp ./src && ./debugApp

update:
  go get -u
  go mod tidy -v




