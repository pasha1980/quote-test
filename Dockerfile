FROM golang:1.22.1 as builder

WORKDIR /app

RUN apt update && apt install -y ca-certificates && rm -rf /var/cache/apk/*

COPY go.mod .

RUN export GOPRIVATE=git.com && go mod download
COPY . .

RUN mkdir -p /tmp_for_scratch

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/quotes

FROM scratch
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tmp_for_scratch /tmp

ENTRYPOINT ["./main"]