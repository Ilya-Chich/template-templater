FROM golang:1.24.0-alpine3.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api


FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/css ./css
COPY --from=builder /app/js ./js
COPY --from=builder /app/img ./img
COPY --from=builder /app/robots.txt /app/sitemap.xml /app/site.webmanifest ./
COPY --from=builder /app/favicon.ico /app/favicon-16x16.png /app/favicon-32x32.png /app/apple-touch-icon.png /app/android-chrome-192x192.png /app/android-chrome-512x512.png ./
CMD ["/app/main"]
