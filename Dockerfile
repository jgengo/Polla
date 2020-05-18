FROM golang:1.14-stretch as build-env

WORKDIR /go/src/app
ADD . /go/src/app

ARG GO111MODULE=on
ARG CGO_ENABLED=1

RUN go get -d -v ./...

RUN go mod vendor
RUN go build -ldflags "-s -w" -o /go/bin/app cmd/polla/main.go 

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
WORKDIR /db # to create a folder while distroless doesnt have a shell
WORKDIR /
CMD ["/app"]