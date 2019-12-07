# build backend
FROM golang:1.12.5-alpine as build
RUN apk add --update --no-cache ca-certificates git nodejs-npm

WORKDIR /github.com/traPtitech/Jomon

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /Jomon


# build frontend(tentative)

WORKDIR /github.com/traPtitech/Jomon/client
RUN npm ci
RUN npm run build


# run

FROM alpine:3.9
WORKDIR /app

COPY --from=build /Jomon ./

ENTRYPOINT ./Jomon