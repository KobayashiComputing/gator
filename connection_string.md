# Postgres Connection String
postgres://postgres:postgres@localhost:5432/gator

## Example
psql "postgres://postgres:postgres@localhost:5432/gator"

## Migrations...
cd into the sql/schema directory and run:
```
goose postgres <connection_string> down
goose postgres <connection_string> up

# example:
# goose postgres "postgres://wagslane:@localhost:5432/gator" down
```
