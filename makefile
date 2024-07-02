MOCKERY_BIN := $(GOPATH)/bin/mockery
GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
GOTOOL = $(GOCMD) tool

.PHONY: serve tidy test mock

serve:
	$(GOCMD) run cmd/api/main.go

cron-category:
	$(GOCMD) run cmd/cron/main.go -run=category
cron-discounts:
	$(GOCMD) run cmd/cron/main.go -run=discounts
cron-disc-costco:
	$(GOCMD) run cmd/cron/main.go -run=discount-costco
cron:
	$(GOCMD) run cmd/cron/main.go -run=all
	
tidy:
	$(GOMOD) tidy && $(GOMOD) vendor
test:
	$(GOTEST) cmd/test/main.go
mock:
	@echo "Generating mocks for interface $(interface) in directory $(dir)..."
	@$(MOCKERY_BIN) --name=$(interface) --dir=$(dir) --output=./internal/mocks
	cd ./internal/mocks && \
	mv $(interface).go $(filename).go
mig-up:
	$(GOCMD) run cmd/migration/main.go -up
mig-down:
	$(GOCMD) run cmd/migration/main.go -down
coverage:
	$(GOTEST) -v ./...
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOTOOL) cover -html=coverage.out -o coverage.html
seed:
	go run cmd/seed/main.go
download:
	@echo Download go.mod dependencies
	@go mod download
install-tools: download
	@echo Installing tools from cmd/tools/main.go
	@cat cmd/tools/main.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
