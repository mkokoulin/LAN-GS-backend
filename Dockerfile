# syntax=docker/dockerfile:1

# Build

FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

# download Go modules and dependencies
RUN go mod download

COPY . ./

RUN go build -o /taplink

COPY --from=build /taplink /usr/local/bin/taplink

# EXPOSE 8080

# ENTRYPOINT /usr/local/bin/taplink --port 8080

FROM nginx:alpine

# Copy your web application files into the image
COPY /usr/local/bin/taplink /usr/share/nginx/html

# Copy the SSL certificate files into the image
COPY certificate.crt /etc/nginx/certs/certificate.crt
COPY certificate.key /etc/nginx/certs/certificate.key

# Configure Nginx to use the SSL certificate
RUN echo "server { \
    listen 443 ssl; \
    server_name taplinkwidget-mikhail.b4a.run; \
    ssl_certificate /etc/nginx/certs/certificate.crt; \
    ssl_certificate_key /etc/nginx/certs/certificate.key; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
    } \
}" > /etc/nginx/conf.d/default.conf

# Deploy 

FROM debian:latest

WORKDIR /

# Expose the HTTPS port
EXPOSE 443