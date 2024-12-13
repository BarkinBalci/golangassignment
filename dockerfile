FROM golang:1.23.4-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o /bin/api ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/api /bin/api

EXPOSE 8080

CMD ["/bin/api"]