FROM golang

ADD . /redcoins

RUN go get golang.org/x/net/html
RUN go get github.com/gorilla/mux
RUN go get github.com/jinzhu/gorm
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/joho/godotenv 

RUN go build

ENTRYPOINT /redcoins

EXPOSE 8000