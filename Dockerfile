FROM golang:1.22-alpine AS builder

WORKDIR /social-app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -mod=vendor -ldflags="-w -s" -o /social-app/social_binary ./cmd/social

FROM alpine:3.18

WORKDIR /social-app

COPY --from=builder /social-app/social_binary .

CMD ["sleep", "infinity"]