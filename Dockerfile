# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

ENV SERVER_ADDRESS=:8002
ENV SECRET_PATH=lan-site-94255-b56432cf737d.json
ENV SCOPE=https://www.googleapis.com/auth/spreadsheets
ENV EVENT_SPREADSHEET_ID=1zssMHkizrIetXEMkV3Qo-wj6QBiv9jf3A2S-g5IaoE0
ENV EVENT_READ_RANGE=master!2:1000
ENV ENTRIES_SPREADSHEET_ID=1IglEBmeCFs9FwL0bQ1vh93b-wJp6KTWurRX2sxcsd3A

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /taplink-backend

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/taplink-backend"]