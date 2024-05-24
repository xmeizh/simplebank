FROM --platform=$BUILDPLATFORM golang:1.22-alpine3.20 as buildstage
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

FROM --platform=$BUILDPLATFORM alpine:3.20
WORKDIR /app
COPY --from=buildstage /app/main .
COPY --from=buildstage /app/migrate ./migrate
COPY app.env .
COPY entrypoint.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/entrypoint.sh" ]
