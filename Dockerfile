# Start from the official Go base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY backend/go.mod backend/go.sum ./backend/

# Download the Go module dependencies
RUN cd backend && go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN cd backend && go build -o click-wars
# RUN cp ./backend/click-wars /app
# RUN cp ./frontend /app

# Expose port 8080 to the outside
EXPOSE 8080

# Run the Go application
CMD ["cd backend && ./click-wars"]