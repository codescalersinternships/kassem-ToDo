FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
ENV DB_FILE="./todo.db"
ENV Port=":5000"
RUN go build -o main -a -ldflags "-linkmode external -extldflags '-static' -s -w" .
FROM alpine:3.11.3
COPY --from=builder /app/main .
EXPOSE 5000

CMD [ "./main" ]