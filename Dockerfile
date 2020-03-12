FROM golang:alpine

WORKDIR /app

ENV API_PORT  8080
ENV CASSANDRA_HOST  localhost
ENV CASSANDRA_USERNAME cassandra
ENV CASSANDRA_PASSWORD cassandra

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Export necessary port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
