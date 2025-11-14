FROM golang:1.22.3-alpine AS dev
WORKDIR /app
COPY . .
RUN apk update && apk upgrade && go mod download
EXPOSE 8000
CMD ["go", "run", "server.go"]
FROM golang:1.22.3-alpine AS app
WORKDIR /app
COPY . .
RUN apk update && apk upgrade && go mod download
EXPOSE 8000
CMD ["go", "run", "server.go"]
