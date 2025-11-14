FROM golang:1.24-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o auth-api server.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/auth-api ./auth-api
COPY certs/ certs/
EXPOSE 8000
CMD ["./auth-api"]
