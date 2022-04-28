# reddit-feed

## Docker-compose
Runs mongo-db and mongo-express
```sh
docker-compose up -d
```

## Migrations

#### Install [Golang-migrate]
```shell
brew install golang-migrate
```

#### Run migrations

```sh
make migrate-up
```

#### Run migrations-down

```sh
make migrate-down
```

## Run tests

Run all tests except integration file (excluded):

```sh
make run-tests
```

Run integration tests:

```sh
make run-integration-tests
```

## Run the app

```shell
make run
```

### Notes

- Remove test data from the cache - $go clean -testcache

[Golang-migrate]: <https://github.com/golang-migrate/migrate>