.PHONY: setup build test lint vet gosec vuln check

setup:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

build:
	go build ./...

test:
	go test ./... -race -count=1 -cover

lint:
	golangci-lint run

vet:
	go vet ./...

gosec:
	gosec ./...

vuln:
	govulncheck ./...

check: vet lint test gosec vuln
