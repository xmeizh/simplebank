# Simple Bank
This application is built during this [Backend Master Class Course.](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/)

## Getting started
### Prerequisites
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://sqlc.dev/)
- [DBDiagram](https://www.dbdiagram.io/d/Simple-bank-66221b7303593b6b6167e52a)

### Usage
#### Initialize DB
```bash
$ make initdb
```

#### Migrate DB Schema
Migrate to an older/newer version:
```bash
$ make [migratedown|migrateup]
```

#### Run Server
```bash
$ make server
```

#### Run Test
```bash
$ make test
```
