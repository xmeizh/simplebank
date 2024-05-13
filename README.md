# Simple Bank
## Tools
- golang-migrate
- [SQLC](https://sqlc.dev/)
- [DB diagram](https://www.dbdiagram.io/d/Simple-bank-66221b7303593b6b6167e52a)

## Golang Dependencies
- [pq](https://github.com/lib/pq)
- [testify](https://github.com/stretchr/testify)

## Notes
### Generate CRUD Golang code from SQL
1. db/sql: fast and straightforward but easy to make mistakes
2. GORM: functions already implemented, slow on high load
3. SQLX: quite fast & easy to use; fields mapping via query text & struct tags; failure won't occur until runtime
4. SQLC: very fast & easy to use; automatic code generation; catch SQL query errors before generating codes; Full support Postgres, MySQL is experimental

## References
- [Backend Master Class (Golang + Postgres + Kubernetes + gRPC)](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/)