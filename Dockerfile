FROM golang:alpine as build-env

# Set the working directory
WORKDIR /go/src/app

# Copy the source code into the container
COPY . .

# Install build-dependencies
RUN apk add --no-cache git

# Build the Go application
RUN go build -o catuaba .

WORKDIR /home
ENTRYPOINT [ "/go/src/app/catuaba" ]
