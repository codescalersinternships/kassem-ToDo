
FROM ubuntu:20.04
WORKDIR /app
RUN apt-get update
ENV DEBIAN_FRONTEND=noninteractive \
TZ=Etc/UTC 
RUN apt install -y golang-go openssh-server nginx
COPY default.conf /etc/nginx/conf.d/
COPY . /app
RUN go mod download
ENV DB_FILE="./todo.db"
ENV Port=":5000"
RUN go build -o main
EXPOSE 5000
RUN chmod +x start.sh
CMD ["/app/start.sh"]
