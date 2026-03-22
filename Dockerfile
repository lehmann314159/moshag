FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o moshag ./cmd/server/main.go

FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/moshag .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
ENV PORT=3000
EXPOSE 3000
CMD ["./moshag"]
