## build backend
FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS server-build

WORKDIR /github.com/traPtitech/Jomon
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /Jomon -ldflags "-s -w"

## run

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
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
