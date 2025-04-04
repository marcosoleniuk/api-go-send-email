FROM golang:1.24-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/main main.go

FROM alpine:latest

RUN apk update && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime && \
    echo "America/Sao_Paulo" > /etc/timezone && \
    apk del tzdata

RUN mkdir /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /.env

EXPOSE 8080

CMD ["/app/main"]
