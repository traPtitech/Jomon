FROM golang:1.16.0-alpine as build
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

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN go get -u github.com/cosmtrek/air

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
