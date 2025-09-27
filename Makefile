.PHONY: install-dev
install-dev:
	@echo "Installing..."
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

.PHONY: fmt
fmt:
	@echo "Formatting..."
	$(shell go env GOPATH)/bin/gci write -s standard -s default -s "prefix(github.com/cugu/md)" .
	$(shell go env GOPATH)/bin/gofumpt -l -w .
	@echo "Done."

.PHONY: lint
lint:
	@echo "Linting..."
	golangci-lint run
	@echo "Done."

.PHONY: test
test:
	@echo "Testing..."
	go test -v ./...

.PHONY: coverage
coverage:
	@echo "Testing with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	@echo "Done."