FROM hub.hamdocker.ir/golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
EXPOSE 8000
ENTRYPOINT ["go", "run", "./cmd/server/main.go"]

