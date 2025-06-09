FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go-api /app/main.go

FROM debian:stable-slim as runtime

WORKDIR /app

COPY --from=builder /app/go-api .

EXPOSE 3000

CMD ["./go-api"]