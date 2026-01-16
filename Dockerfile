FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build publisher
RUN CGO_ENABLED=0 go build -o /publisher ./cmd/publisher

# Build subscriber
RUN CGO_ENABLED=0 go build -o /subscriber ./cmd/subscriber

# Publisher runtime
FROM alpine:latest AS publisher
RUN apk --no-cache add ca-certificates
COPY --from=builder /publisher /publisher
ENTRYPOINT ["/publisher"]

# Subscriber runtime
FROM alpine:latest AS subscriber
RUN apk --no-cache add ca-certificates
COPY --from=builder /subscriber /subscriber
ENTRYPOINT ["/subscriber"]
