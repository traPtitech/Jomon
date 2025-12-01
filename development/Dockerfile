FROM golang:alpine as build

ENV CGO_ENABLED 0
ENV TZ Asia/Tokyo



RUN apk add --update --no-cache git
RUN apk --update --no-cache add tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && apk del tzdata
RUN apk --update --no-cache add ca-certificates \
    && update-ca-certificates \
    && rm -rf /usr/share/ca-certificates /etc/ssl/certificates

RUN go install github.com/powerman/dockerize@v0.24.0
RUN go install github.com/go-task/task/v3/cmd/task@v3.45.5

RUN mkdir -p /app

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
