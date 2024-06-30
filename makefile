MOCKERY_BIN := $(GOPATH)/bin/mockery

.PHONY: serve tidy test mock

serve:
	go run cmd/api/main.go
tidy:
	go mod tidy && go mod vendor
test:
	go run cmd/test/main.go
mock:
	@echo "Generating mocks for interface $(interface) in directory $(dir)..."
	@$(MOCKERY_BIN) --name=$(interface) --dir=$(dir) --output=./internal/mocks
	cd ./internal/mocks && \
	mv $(interface).go $(filename).go
mig-up:
	go run cmd/migration/main.go -up
mig-down:
	go run cmd/migration/main.go -down
coverage:
	go test -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
seed:
	go run cmd/seed/main.go
download:
	@echo Download go.mod dependencies
	@go mod download
install-tools: download
	@echo Installing tools from cmd/tools/main.go
	@cat cmd/tools/main.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
