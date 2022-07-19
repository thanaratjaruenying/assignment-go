# syntax=docker/dockerfile:1

# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.18-alpine

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod go.sum /app/

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . /app/

# compile application
RUN cd /app
RUN go build -o /app/healthcheck
 
# command to be used to execute when the image is used to start a container
CMD [ "/app/healthcheck" ]
