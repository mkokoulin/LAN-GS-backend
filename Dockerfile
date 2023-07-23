# syntax=docker/dockerfile:1

# Build
FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY cert.pem .
COPY key.pem .

# download Go modules and dependencies
RUN go mod download
COPY . ./

RUN go build -o /taplink

# Deploy 

FROM debian:latest
WORKDIR /
COPY --from=build /taplink /usr/local/bin/taplink
EXPOSE 8080
ENTRYPOINT /usr/local/bin/taplink --port 8080