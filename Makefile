build:
	go build -o bin/crash_reporter

unit_test:
	go test -v ./cmd/...

e2e_test: build e2e_test_runner

e2e_test_local: build e2e_local_setup e2e_test_runner

e2e_test_runner:
	go test -v ./e2e/...

e2e_local_setup:
	cp test_credentials/gcp_credentials.json /tmp/
	cp e2e/helpers/test_config.yaml /tmp/