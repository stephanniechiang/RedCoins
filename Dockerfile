FROM golang

ADD . /go/src/redcoins

RUN go get golang.org/x/net/html
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux

RUN go build

ENTRYPOINT /go/bin/redcoins

EXPOSE 8080