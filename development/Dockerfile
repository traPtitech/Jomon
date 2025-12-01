FROM golang:alpine as build
ENV DOCKERIZE_VERSION v0.6.1
ENV CGO_ENABLED 0
ENV TZ Asia/Tokyo

ARG GO_BUILD_TAGS=''
RUN echo $GO_BUILD_TAGS

RUN apk add --update --no-cache git
RUN apk --update --no-cache add tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && apk del tzdata
RUN apk --update --no-cache add ca-certificates \
    && update-ca-certificates \
    && rm -rf /usr/share/ca-certificates /etc/ssl/certificates

RUN go install github.com/powerman/dockerize@latest
RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
