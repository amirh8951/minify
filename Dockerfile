FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /build/url-shortener ./cmd/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/url-shortener /url-shortener

EXPOSE 5000

ENTRYPOINT ["/url-shortener"]
