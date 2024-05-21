# Simple Bank
This application is built during this [Backend Master Class Course.](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/)

## Usage
### Initialize DB
```bash
$ make initdb
```

### Migrate DB Schema
Migrate to an older/newer version:
```bash
$ make [migratedown1|migrateup1]
```

### Run Server
```bash
$ make server
```

### Run Test
```bash
$ make test
```

## Tools
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://sqlc.dev/)
- [DBDiagram](https://www.dbdiagram.io/d/Simple-bank-66221b7303593b6b6167e52a)
- [gomock](https://github.com/golang/mock?tab=readme-ov-file)
