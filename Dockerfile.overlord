# Use an official Golang runtime as a parent image
FROM golang:alpine3.17 AS build

# Set the working directory to the app directory
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies. They will be cached if the go.mod and go.sum files don't change
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the binary executable
RUN go build -o myapp .

# Start a new stage from scratch to create a smaller final image
FROM golang:alpine3.17

# Set the working directory to the app directory
WORKDIR /app

# Copy the binary executable from the previous stage
COPY --from=build /app/myapp ./

EXPOSE 3000

# Run the binary executable
CMD ["./myapp"]
