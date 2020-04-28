## build backend
FROM golang:1.13.5-alpine as build
RUN apk add --update --no-cache git nodejs-npm

WORKDIR /github.com/traPtitech/Jomon
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /Jomon

## build frontend
WORKDIR /github.com/traPtitech/Jomon/client
COPY ./client .
RUN npm ci
RUN npm run build

## run

FROM alpine:3.9
WORKDIR /app
RUN apk --update add tzdata ca-certificates && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  rm -rf /var/cache/apk/*
COPY --from=build /Jomon ./
COPY --from=build /github.com/traPtitech/Jomon/client/dist ./client/dist/
ENTRYPOINT ./Jomon
 