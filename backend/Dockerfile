
FROM ubuntu:20.04
WORKDIR /app
RUN apt-get update
ENV DEBIAN_FRONTEND=noninteractive \
TZ=Etc/UTC 
RUN apt install -y golang-go openssh-server nginx
COPY default.config /etc/nginx/conf.d/
COPY . /app
RUN go mod download
ENV DB_FILE="./todo.db"
ENV Port=":5000"
RUN go build -o main
EXPOSE 5000
RUN chmod +x start.sh
CMD ["/app/start.sh"]
# FROM golang:latest as builder

# WORKDIR /app

# COPY go.mod /app
# COPY go.sum /app
# RUN go mod download

# COPY . .
# ENV DB_FILE="./todo.db"
# ENV Port=":5000"
# RUN go build -o main -a -ldflags "-linkmode external -extldflags '-static' -s -w" .
# FROM alpine:3.11.3
# COPY --from=builder /app/main /app
# EXPOSE 5000

# CMD [ "./main" ]