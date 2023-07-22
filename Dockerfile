# # syntax=docker/dockerfile:1

# # specify the base image to  be used for the application, alpine or ubuntu
# FROM golang:1.19.5-alpine

# # create a working directory inside the image
# WORKDIR /app

# ENV SERVER_ADDRESS=:8002
# ENV SCOPE=value1
# ENV EVENT_SPREADSHEET_ID=value2
# ENV EVENT_READ_RANGE=value3
# ENV ENTRIES_SPREADSHEET_ID=value4

# # copy Go modules and dependencies to image
# COPY go.mod go.sum ./

# # download Go modules and dependencies
# RUN go mod download

# # copy directory files i.e all files ending with .go
# COPY *.go ./

# # compile application
# RUN go build -o /dist

# ##
# ## STEP 2 - DEPLOY
# ##
# FROM scratch

# WORKDIR /

# COPY --from=build /dist /dist

# # tells Docker that the container listens on specified network ports at runtime
# EXPOSE 80
# EXPOSE 8080
# EXPOSE 3306
# EXPOSE 27018

# ENTRYPOINT ["/dist"]

# syntax=docker/dockerfile:1

# Build

FROM golang:1.18 AS build

WORKDIR /app

ENV SCOPE=value0
ENV EVENT_SPREADSHEET_ID=value1
ENV EVENT_READ_RANGE=value2
ENV ENTRIES_SPREADSHEET_ID=value3
ENV GOOGLE_SECRET=value4

COPY go.mod .
COPY go.sum .

# download Go modules and dependencies
RUN go mod download

COPY . ./

RUN go build -o /taplink

# Deploy 

FROM debian:latest

WORKDIR /

COPY --from=build /taplink /usr/local/bin/taplink

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/taplink" ]