# logger-api

## Setup

- Set `LOGGER_STATIC_FILES` environment variable to the `client/` path
- Make a copy of `dev.yml.example` and save it as dev.yml
  - Update `connection_string` and `database_name` eg:
    ```
    connection_string: root:@tcp(127.0.0.1:2020)/
    database: logger_dev
    ```
- Do the same for `test.yml.example`
- Run tests
  ```
  $ go test $(go list ./... | grep -v /vendor/)
  ```
- Run migrations in `dev` `prod`: todo
