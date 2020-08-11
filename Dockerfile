FROM golang:alpine

WORKDIR /app

ENV MYSQL_DRIVER mysql
ENV MYSQL_USER root
ENV MYSQL_PASSWORD admin
ENV MYSQL_ROOT_PASSWORD admin
ENV MYSQL_DATABASE api_go

ADD database/schema-apigo.sql /docker-entrypoint-initdb.d

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Export necessary port
EXPOSE 8000:8000

# Command to run the executable
CMD ["./main"]
