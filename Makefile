APP?=gotest2allure
RELEASE?=1.0.0
GOOS?=darwin

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build \
		-ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o bin/${GOOS}/${RELEASE}/${APP}

.PHONY: clean
clean:
	@rm -f bin/${GOOS}/${APP}