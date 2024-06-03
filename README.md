# Simple Bank
This application is built during this [Backend Master Class Course.](https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/)

## Usage
### Run the application on a local machine:
```bash
$ docker compose up
```

### Run the unit tests:
```bash
$ make test
```

## Tools
- [golang-migrate:](https://github.com/golang-migrate/migrate) helps migrate database schemata.
- [sqlc:](https://sqlc.dev/) generates sql codes in Go.
- [DBDiagram:](https://www.dbdiagram.io/d/Simple-bank-66221b7303593b6b6167e52a) is a great tool to draw a database schema and export it as a sql query.
- [gomock:](https://github.com/golang/mock?tab=readme-ov-file) is used for mock tests.
- [cert-manager:](https://cert-manager.io/) renews TLS certificates.

## Documentation
* [Configuring OpenID Connect in AWS](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services)
* [Grant access to Kubernetes APIs](https://docs.aws.amazon.com/eks/latest/userguide/grant-k8s-access.html)