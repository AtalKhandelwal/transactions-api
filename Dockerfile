# --- builder ---
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Cache deps
COPY go.mod ./
RUN go mod download

# Copy source then ensure go.sum is generated
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server

# --- runtime ---
FROM alpine:3.20
RUN adduser -D -H app && apk add --no-cache ca-certificates
USER app
WORKDIR /home/app
COPY --from=builder /server ./server
EXPOSE 8080
ENV PORT=8080
CMD ["./server"]
