FROM golang:alpine

MAINTAINER David Budworth <dbudworth@peak6.com>

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["app"]

EXPOSE 8080

COPY . /go/src/app
RUN go install && rm -r /go/pkg
