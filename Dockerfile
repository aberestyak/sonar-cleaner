ARG GOLANG_VERSION=1.18.0-alpine3.15
ARG IMAGE_TAG__ALPINE_GO=3.15

FROM golang:${GOLANG_VERSION} as build
WORKDIR /go/src/sonar-cleaner
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go generate ./...
ENV GOGCFLAGS= \
    CGO_ENABLED=0 \
    GO111MODULE=on

ARG ARTIFACT_VERSION=local
RUN go install \
    -installsuffix "static" \
    -ldflags "                                            \
        -X main.Version=${ARTIFACT_VERSION}               \
        -X main.GoVersion=$(go version | cut -d " " -f 3) \
        -X main.Compiler=$(go env CC)                     \
        -X main.Platform=$(go env GOOS)/$(go env GOARCH)  \
    " \
    ./...


FROM alpine:${IMAGE_TAG__ALPINE_GO} as runtime
RUN apk --no-cache upgrade \
    && apk --no-cache add tzdata \
    && echo 'Etc/UTC' > /etc/timezone
ENV TZ     :/etc/localtime \
    LANG   en_US.utf8 \
    LC_ALL en_US.UTF-8
COPY --from=build /go/bin/sonar-cleaner /sonar-cleaner
ENTRYPOINT ["/sonar-cleaner"]
