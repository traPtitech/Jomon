FROM golang:1.13.5-alpine as build
ENV DOCKERIZE_VERSION v0.6.1
ENV CGO_ENABLED 0
ARG GO_BUILD_TAGS=''

RUN echo $GO_BUILD_TAGS

RUN apk add --update --no-cache ca-certificates git openssl \
    && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download \
    && go get github.com/pilu/fresh

COPY . .

RUN go build -tags $GO_BUILD_TAGS -o /app/Jomon
