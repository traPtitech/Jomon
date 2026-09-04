## build backend
FROM golang:1.27.1-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS server-build

WORKDIR /github.com/traPtitech/Jomon
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /Jomon -ldflags "-s -w"

## run

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
ENV TZ Asia/Tokyo

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates

WORKDIR /app
COPY --from=server-build /Jomon ./
COPY --from=server-build /github.com/traPtitech/Jomon/migrations ./migrations

ENTRYPOINT ./Jomon
