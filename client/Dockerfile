FROM node:lts-alpine AS build
WORKDIR /app

ENV TZ Asia/Tokyo

RUN apk --update --no-cache add git
RUN apk --update --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
