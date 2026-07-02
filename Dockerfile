FROM golang:1.24-alpine AS builder

WORKDIR /app

# Go source files now live under src/ for clean project structure
COPY src/go.mod ./
COPY src/ .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o bot-binary main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bot-binary .

RUN apk add --no-cache tzdata ca-certificates su-exec

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Run natively using Alpine directly without massive OS overhead!
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./bot-binary"]
