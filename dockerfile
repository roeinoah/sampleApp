# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Roei Noah"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
#COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
COPY .aws ~/
# Build the Go app
RUN go get github.com/aws/aws-sdk-go/ 
RUN go build -o Main .

# Expose port 9090 to the outside world
EXPOSE 9090

# Command to run the executable
CMD ["./Main"]
