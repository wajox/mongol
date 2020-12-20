start-docker-compose-test:
	docker-compose -f docker-compose-test.yml up -d

stop-docker-compose-test:
	docker-compose -f docker-compose-test.yml down

test-all:
	$(MAKE) start-docker-compose-test
	go test -v ./...
	${MAKE} stop-docker-compose-test
install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.33.0
	go get github.com/google/wire/cmd/wire
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/swaggo/swag/cmd/swag
