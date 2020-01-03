
#build stage
FROM golang:alpine AS builder
ENV GO111MODULE=on CGO_ENABLED=0
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build cmd/mibolsillo/main.go
#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/app .
RUN  chmod +x main
LABEL Name=mibolsillo Version=0.0.1
EXPOSE 5000
CMD ["./main"]

