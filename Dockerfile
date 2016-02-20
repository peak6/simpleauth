FROM golang:alpine

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["app"]

COPY . /go/src/app
RUN go install 
