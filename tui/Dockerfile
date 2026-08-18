FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tui ./cmd/tui

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/tui .
ENV SSH_LISTEN=:2222
ENV SSH_HOST_KEY=/data/ssh/id_ed25519
EXPOSE 2222
CMD ["./tui"]
