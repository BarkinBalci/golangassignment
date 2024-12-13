FROM golang:1.23.4-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o /app/api ./cmd/api/main.go

FROM golang:1.23.4-bookworm

WORKDIR /app

COPY --from=builder /app/api /app/api

EXPOSE 8080

CMD ["/app/api"]