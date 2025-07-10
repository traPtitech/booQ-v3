#!/bin/sh
docker compose -f docker/test/docker-compose.yml up -d --wait

export MYSQL_DATABASE=booq-test
export MYSQL_PORT=3307

if [ "$1" = "cover" ]; then
  echo "Running Tests With Coverage..."

  go test . ./model -v -covermode=atomic -cover -coverprofile=cover_model.out
  go test . ./router -v -covermode=atomic -cover -coverprofile=cover_router.out

  go tool cover -html=cover_model.out -o cover_model.html
  go tool cover -html=cover_router.out -o cover_router.html
else
  echo "Running Tests..."

  go test . ./model -v -covermode=atomic
  go test . ./router -v -covermode=atomic
fi

docker compose -f docker/test/docker-compose.yml down