FROM golang:1.11


WORKDIR /go/src/go_app
COPY . .

RUN apt-get update 
RUN apt-get install -y vim
ENV GO111MODULE on
RUN go get -d -v ./...
RUN go install -v ./...
RUN go get -u gopkg.in/godo.v2/cmd/godo

CMD ["make", "watch-run"]