# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /otelcol-genai-safe ./cmd/otelcol-custom

# ---------- Runtime stage ----------
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /otelcol-genai-safe /otelcol-genai-safe

ENTRYPOINT ["/otelcol-genai-safe"]
