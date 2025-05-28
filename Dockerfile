# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Final stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/main /app/main

EXPOSE 5555

CMD ["/app/main"]