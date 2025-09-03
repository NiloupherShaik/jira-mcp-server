version := v0.0.1
SERVICE_NAME ?=

all: build

build:
	@CGO_ENABLED=0 GOOS=linux go build -a --ldflags "-X 'main.Version=${version}'-extldflags '-static'" \
		-tags "netgo" -installsuffix netgo \
		-o /go/bin/service \
		cmd/${SERVICE_NAME}/main.go

	go version
