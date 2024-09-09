FROM golang:1.23.1-alpine3.20 AS build
COPY src /tmp/src
WORKDIR /tmp/src
RUN apk add make && make build-migrate && make build

FROM alpine:3.20
EXPOSE 8080
WORKDIR /app
COPY --from=build /tmp/src/bin/app /app/avito_service
COPY --from=build /tmp/src/bin/migrate /app/migrate
COPY --from=build /tmp/src/migrate/migrations /app/migrations
ENTRYPOINT [ "/bin/sh", "-c", "/app/migrate && /app/avito_service" ]
