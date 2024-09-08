# Import Base Image Golang
FROM golang:alpine as compiler

# Setting Working Directory
WORKDIR /app

# Import Source Code
COPY . . 

# Install Depencies and Get Package
RUN go mod download && go mod verify

# Compile Source Code
RUN go build -o /app/main

# Multi Stage Build Docker
FROM alpine:latest

# Setting Working Directory
WORKDIR /app

# Setting Env Variable
ENV APP_MODE=production
ENV HOST=0.0.0.0
ENV PORT=3000

ENV DB_HOST=127.0.0.1
ENV DB_PORT=3306
ENV DB_USERNAME=root
ENV DB_PASSWORD=
ENV DB_NAME=

ENV JWT_SECRET=bwa
ENV JWT_EXPIRED=5

# Copy Compile Golang APP
COPY --from=compiler /app/main /app/main

# Expose
EXPOSE 3000

# Running App
CMD ["/app/main"]