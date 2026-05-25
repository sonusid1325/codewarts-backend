FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency configs
COPY go.mod ./
# If go.sum exists, copy it (might not exist yet, we will create it)
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /

COPY --from=builder /server /server

EXPOSE 8080

CMD ["/server"]
