FROM --platform=$BUILDPLATFORM golang:1.22-alpine3.20 as buildstage
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=buildstage /app/main .
COPY app.env .
COPY entrypoint.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/entrypoint.sh" ]
