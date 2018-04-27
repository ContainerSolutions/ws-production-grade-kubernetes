FROM golang:1.9

# Build and test deps
RUN go get github.com/jinzhu/gorm
RUN go get github.com/jinzhu/gorm/dialects/postgres
RUN go get github.com/jinzhu/gorm/dialects/sqlite

COPY *.go /go/
RUN go build -o /home/server
COPY test.sh /
RUN chmod +x /test.sh

ENTRYPOINT ["/home/server"]

EXPOSE 8080
ENV VERSION 2.0
