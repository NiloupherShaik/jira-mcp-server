# Step 1 - Build executable
FROM public.ecr.aws/docker/library/golang:1.24.6-bullseye AS build-stage

WORKDIR /go/src/app

RUN apt-get update &&\
    apt-get install -y git &&\
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

ARG GITHUB_USERNAME
ARG GITHUB_PASSWORD
ARG ARTIFACTORY_USERNAME
ARG ARTIFACTORY_PASSWORD

RUN git config --global credential.username $GITHUB_USERNAME &&\
    git config --global credential.helper '!f() { echo "password=${GITHUB_PASSWORD}"; }; f' &&\
    go env -w GOPRIVATE=github.com/nable-fusion/*,github.com/logicnow/* &&\
    go env -w GONOPROXY=NONE &&\
    go env -w GOPROXY="https://${ARTIFACTORY_USERNAME}:${ARTIFACTORY_PASSWORD}@mspsolarwinds.jfrog.io/artifactory/api/go/davinci-go,https://proxy.golang.org,direct" &&\
    go mod download

COPY internal/pkg ./internal/pkg
COPY build/package/Makefile ./Makefile

ARG SERVICE_NAME
COPY cmd/${SERVICE_NAME}/ ./cmd/${SERVICE_NAME}/
COPY internal/app/${SERVICE_NAME}/ ./internal/app/${SERVICE_NAME}/
RUN make

# Step 2 - Run application
FROM public.ecr.aws/docker/library/alpine:3.20.0 as release-stage
COPY --from=build-stage /go/bin/service /service

CMD ["/service"]
