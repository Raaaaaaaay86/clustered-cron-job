# Build
FROM golang:1.19.2-buster AS build

COPY . /gin
WORKDIR /gin

RUN go mod download
RUN go build -o gin-server

# Deploy
FROM ubuntu:latest

WORKDIR /

COPY --from=build /gin/.env /.env
COPY --from=build /gin/gin-server /gin

ENTRYPOINT ["/gin"]
