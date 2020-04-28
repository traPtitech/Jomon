## build backend
FROM golang:1.13.5-alpine as server-build

WORKDIR /github.com/traPtitech/Jomon
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /Jomon

## build frontend
FROM node:13.12.0-alpine as client-build
WORKDIR /github.com/traPtitech/Jomon/client
COPY ./client .
RUN npm ci
RUN npm run build

## run

FROM alpine:3.9
ENV DOCKERIZE_VERSION v0.6.1
RUN apk add --update --no-cache ca-certificates git openssl curl \
    && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
WORKDIR /app
RUN apk --update add tzdata ca-certificates && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  rm -rf /var/cache/apk/*
COPY --from=server-build /Jomon ./
COPY --from=client-build /github.com/traPtitech/Jomon/client/dist ./client/dist/
ENTRYPOINT ./Jomon
 