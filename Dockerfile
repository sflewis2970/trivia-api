# Base the image off of the lastest version of 1.18
FROM golang:1.18

# Make the detsination directory
RUN mkdir -p /home/app

# Copy files to the new directory
COPY . /home/app

# Set working directory so that the build command can locate the go.mod file
WORKDIR /home/app

# Build the application
RUN go build -v -o ./main ./cmd/services

# Run the app in the image
CMD ["./main"]
