docker compose -f docker/test/docker-compose.yml up -d

export MYSQL_DATABASE=booq-test
export MYSQL_PORT=3307

go test . ./model -v -covermode=atomic -vet=off
go test . ./router -v -covermode=atomic -vet=off

docker compose -f docker/test/docker-compose.yml down