# Import Go image. The application code was written in Go version 1.15.6
FROM golang:1.15.6

# Set app directory
WORKDIR $GOPATH/src/notetaker-app

# Copy the application code from the local machine to the working directory
COPY . .

# Download required packages and install the application
RUN go get -d -v ./...
RUN go install -v ./...

# expose the HTTP port on our app (80)
EXPOSE 80

# start app
CMD ["golang-notetaker"]