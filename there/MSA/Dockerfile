FROM golang:1.13

WORKDIR /go/src/app
COPY . .

EXPOSE 8888

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
