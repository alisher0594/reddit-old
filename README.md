![Reddit-old](docs/img/logo.png)

#  Old Reddit 

a simplified version of the Reddit feed API

## Content
- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [Dependency Injection](#dependency-injection)


## Quick start
Local development:
```sh
# Postgres
$ make compose-up
# Run app with migrations
$ make run
```

Integration tests (can be run in CI):
```sh
# DB, app + migrations, integration tests
$ make compose-up-integration-test
```

## Project structure
### `cmd/api`
Configuration and logger initialization.
This is where all the main objects are created.
Dependency injection occurs through the "New ..." constructors (see Dependency Injection).
This technique allows us to layer the application using the [Dependency Injection](#dependency-injection) principle.
This makes the business logic independent from other layers.

Next, we start the server and wait for signals for graceful completion.

The `migrate.go` file is used for database auto migrations.
It is included if an argument with the _migrate_ tag is specified.
For example:

```sh
$ go run -tags migrate ./cmd/api
```

### `integration-test`
Integration tests.
They are launched as a separate container, next to the application container.
It is convenient to test the Rest API using [go-hit](https://github.com/Eun/go-hit).

### `internal/data/entitys`
Entities of business logic (models) can be used in any layer.
There are also methods, for validation.


## Dependency Injection
In order to remove the dependence of business logic on external packages, dependency injection is used.

For example, through the New constructor, we inject the dependency into the structure of the business logic.
This makes the business logic independent (and portable).