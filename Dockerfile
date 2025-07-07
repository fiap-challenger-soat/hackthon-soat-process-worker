FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc g++ make git ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker/main.go

FROM alpine:latest

RUN apk add --no-cache ffmpeg

WORKDIR /root/

COPY --from=builder /app/worker .

CMD ["./worker"]