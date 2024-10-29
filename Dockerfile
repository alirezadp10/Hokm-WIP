FROM hub.hamdocker.ir/golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -o app
#
#FROM alpine:latest
#WORKDIR /app
#COPY --from=builder /app/app .
#ENTRYPOINT ["./app","serve"]

ENTRYPOINT ["go", "run", "main.go"]
