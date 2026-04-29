FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /bin/devops ./

FROM scratch

COPY --from=builder /bin/devops /devops
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/devops"]
