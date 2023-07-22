# syntax=docker/dockerfile:1

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.19-alpine

# create a working directory inside the image
WORKDIR /app

ENV SERVER_ADDRESS=:8002
ENV SCOPE=value1
ENV EVENT_SPREADSHEET_ID=value1
ENV EVENT_READ_RANGE=value1
ENV ENTRIES_SPREADSHEET_ID=value1

# copy Go modules and dependencies to image
COPY go.mod go.sum ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY *.go ./

# compile application
RUN CGO_ENABLED=0 GOOS=linux go build -o /dist

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 80
EXPOSE 8080
EXPOSE 3306
EXPOSE 27018

# command to be used to execute when the image is used to start a container
CMD [ "/dist" ]