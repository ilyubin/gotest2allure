APP?=gotest2allure
RELEASE?=0.0.1
GOOS?=darwin

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: buildt
buildt: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build \
		-ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o bin/${GOOS}/${RELEASE}/${APP} \
		cmd/gotest2allure/main.go

.PHONY: clean
clean:
	@rm -f bin/${GOOS}/${APP}

.PHONY: test
test:
	go test -v ./...

.PHONY: all
all:
	ls

.PHONY: lint
lint:
	golangci-lint run
