# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN mkdir /go-distributed-cache
RUN go build -o /go-distributed-cache ./...

EXPOSE 8080
ENTRYPOINT [ "/go-distributed-cache/cmd" ]