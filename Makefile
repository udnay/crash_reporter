build:
	go build -o bin/crash_reporter

unit_test:
	go test -v ./cmd/...