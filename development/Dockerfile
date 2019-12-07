FROM golang:1.12.5-alpine as build
RUN apk add --update --no-cache ca-certificates git

WORKDIR /github.com/traPtitech/Jomon

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /Jomon

FROM alpine:3.9
WORKDIR /app

COPY --from=build /Jomon ./

ENTRYPOINT ./Jomon