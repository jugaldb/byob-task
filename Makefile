GO=go
GOVET=$(GO) vet
GORUN=$(GO) run
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover
GOBUILD=$(GO) build
COVTHRESHOLD=70
.PHONY: test migrations migrate showmigrations revertmigrations setmigrations
.SILENT: test
all: check test build
check: ## Run static checks
	$(GOVET) ./...
test: ## Execute test cases with code coverage
	$(GOTEST) -v -race -covermode=atomic -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html
	COVERAGE=$$(go tool cover -func=coverage.out | grep "^total:" | grep -o "[0-9\.]*");\
	echo "$$COVERAGE $(COVTHRESHOLD)" | awk '{if (!($$1 >= $$2)) { print "Coverage below threshold - Coverage: " $$1 "%" ", Expected threshold: " $$2 "%"; exit 1 } else { print "Coverage above threshold - Coverage: " $$1 "%" ", Expected threshold: " $$2 "%"; } }'
build: ## Clean dist directory and rebuild the binary file
	rm -rf ./dist && CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -o ./dist/app ./src