# Use an official Golang runtime as a parent image
FROM golang:1.21

# Setting port
ARG APP_PORT="8080"
ENV APP_PORT=$APP_PORT

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application code to the container's workspace
COPY . .

# Build the Go application
RUN go build -o myapp

# Expose the port the application will run on
EXPOSE $APP_PORT

# Define the command to run your application
CMD ["./myapp"]
