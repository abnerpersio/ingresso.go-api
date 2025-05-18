FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN go build -o go-api /app/api/main.go

FROM debian:stable-slim as runtime

WORKDIR /app

COPY --from=builder /app/go-api .

EXPOSE 3000

CMD ["./go-api"]