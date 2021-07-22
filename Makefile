build:
	go build -o bin/crash_reporter

unit_test:
	go test -v ./cmd/...

e2e_test: build e2e_test_runner

e2e_test_runner:
	go test -v ./e2e/...