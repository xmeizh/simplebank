# Tools
- golang-migrate
- [SQLC](https://sqlc.dev/)

# Usage

# Notes
## Generate CRUD Golang code from SQL
1. db/sql: fast and straightforward but easy to make mistakes
2. GORM: functions already implemented, slow on high load
3. SQLX: quite fast & easy to use; fields mapping via query text & struct tags; failure won't occur until runtime
4. SQLC: very fast & easy to use; automatic code generation; catch SQL query errors before generating codes; Full support Postgres, MySQL is experimental