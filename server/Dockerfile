FROM golang:1.8.3

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

# set up autoreload
RUN go get github.com/handwritingio/autoreload

EXPOSE 8080

#
# CMD ["go-wrapper", "run"] # ["app"]
CMD ["/go/bin/autoreload"]
