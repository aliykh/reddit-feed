# reddit-feed

#### Environment setup - MongoDB and Mongo express

```sh
docker-compose up -d
```

#### Run migrations

```sh
make migrate-up
```

#### Run migrations-down

```sh
make migrate-down
```

#### Run tests

Run all tests except integration file (excluded):

```sh
make run-tests
```

Run integration tests:

```sh
make run-integration-tests
```

### Run the app

```shell
make run
```

### Notes

- Remove test data from the cache - $go clean -testcache