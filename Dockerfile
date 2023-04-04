FROM golang:1.19 AS builder

WORKDIR /app

COPY .. .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/app ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/app /bin/app

ENTRYPOINT ["/bin/app"]
