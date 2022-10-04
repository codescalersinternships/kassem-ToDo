
FROM golang:latest as builder
WORKDIR /app

COPY . /app
ENV DB_FILE="./todo.db"
ENV Port=":5000"
RUN go build -o main -a -ldflags "-linkmode external -extldflags '-static' -s -w" .
FROM alpine:3.11.3
COPY start.sh /
COPY --from=builder /app/main /app

EXPOSE 5000

CMD ["./app/main"]

