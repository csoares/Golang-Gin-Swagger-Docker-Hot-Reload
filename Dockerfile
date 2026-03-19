FROM golang:1.21-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc libc-dev

# Install swag (swagger generator) and CompileDaemon (hot reload)
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3 && \
    go install github.com/githubnemo/CompileDaemon@latest

# Set the Current Working Directory inside the container
WORKDIR /go/src/projetoapi

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy

# RUN Swagger
RUN swag init

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# HOT RELOAD
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build ./main.go" -command="./main"
