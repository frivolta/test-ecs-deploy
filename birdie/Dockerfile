# Build stage
FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /app
COPY .. .
RUN apk update; apk add curl
RUN go build -o main main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
# Replacing env variables
# COPY app.env .
COPY serviceAccountKey.prod.json ./serviceAccountKey.json
COPY app.env.production ./app.env
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migrations
EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]