FROM golang:latest AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/yamlSplitter .
ENTRYPOINT ["/app/yamlSplitter"]