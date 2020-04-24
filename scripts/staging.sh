#!/bin/bash

cd docker/staging

git fetch origin master
diff="git diff HEAD..origin/master --name-only"

docker=`$diff | grep docker/`
client=`$diff | grep client/`
server=`$diff | grep .go`

git pull origin master

if [ -n "$docker" ] ; then
  echo "all deploying..."
  docker-compose build --no-cache
  docker-compose up -d --force-recreate
  docker image prune -f
  exit 0
fi

if [ -n "$server" ] ; then
  echo "server deploying..."
  docker-compose build jomon-server
  docker-compose up -d --force-recreate jomon-server
fi

if [ -n "$client" ] ; then
  echo "client deploying..."
  docker-compose build --no-cache jomon-client
  docker-compose down
  docker-compose up -d
fi

docker image prune -f