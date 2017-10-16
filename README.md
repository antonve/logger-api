# logger-api

[![CircleCI](https://circleci.com/gh/antonve/logger-api/tree/master.svg?style=svg)](https://circleci.com/gh/antonve/logger-api/tree/master)

## Setup
- Install glide: https://github.com/Masterminds/glide
- Install dependencies
  ```
  $ glide install
  ```
- Set `LOGGER_STATIC_FILES` environment variable to the `client/` path once the client is available
- Make a copy of `dev.yml.example` and save it as dev.yml
  - Uses PostgreSQL
  - Update `connection_string` and `database_name` eg:
    ```
    connection_string: user=anton sslmode=disable dbname=
    database: logger_dev
    ```
- Do the same for `test.yml.example`
- Run tests
  ```
  $ go test $(go list ./... | grep -v /vendor/)
  ```
- Run migrations in `dev` `prod`: todo
