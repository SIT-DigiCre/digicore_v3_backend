# digicore v3 backend

## setup env

```sh
go mod download
set -a && source .env && set +a
go run github.com/rubenv/sql-migrate/sql-migrate@v1.1.1 up
```
