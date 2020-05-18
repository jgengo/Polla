FROM golang:1.14 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/app cmd/polla/main.go 

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
CMD ["/app"]