FROM node:lts-alpine AS build
WORKDIR /app
RUN apk add --no-cache git && \
    apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

COPY package*.json ./
RUN npm ci
COPY . .
