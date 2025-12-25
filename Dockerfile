# Stage 1: Build the Go binary
FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: Minimal secure runtime
# FROM scratch
FROM alpine

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/static/ static

EXPOSE 8080
ENTRYPOINT ["./server"]
