FROM golang:1.25.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/google/wire/cmd/wire@latest

COPY . .

RUN wire ./cmd/http

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o isayoga-api ./cmd/http

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/isayoga-api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./isayoga-api"]
