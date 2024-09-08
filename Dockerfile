FROM golang:1.23.1-alpine3.20 AS build
COPY src /tmp/src
WORKDIR /tmp/src
RUN apk add make && make build

FROM alpine:3.20
EXPOSE 8080
RUN mkdir /app
COPY --from=build /tmp/src/build/* /app/avito_service
ENTRYPOINT [ "/app/avito_service" ]
