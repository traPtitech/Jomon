## build backend
FROM golang:1.18.0-alpine as server-build

WORKDIR /github.com/traPtitech/Jomon
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /Jomon -ldflags "-s -w"

## run

FROM alpine:3.15.1
ENV TZ Asia/Tokyo

RUN apk --update --no-cache add tzdata \
  && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
  && apk del tzdata
RUN apk --update --no-cache add ca-certificates \
  && update-ca-certificates \
  && rm -rf /usr/share/ca-certificates /etc/ssl/certs

WORKDIR /app
COPY --from=server-build /Jomon ./

ENTRYPOINT ./Jomon
