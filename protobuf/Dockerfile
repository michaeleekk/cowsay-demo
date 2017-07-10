FROM nanoservice/protobuf-go:latest as protos

COPY *.proto /src/
RUN protoc --go_out=. *.proto

ENTRYPOINT []

FROM golang:1.8.3

WORKDIR /go/src/protot
COPY --from=protos /src/*.pb.go ./
WORKDIR /go/src
COPY interface protot_interface/
WORKDIR /go/src/app
COPY *.go ./

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"]
