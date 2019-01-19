FROM golang:1.11-alpine AS builder
 
RUN apk add bash git gcc g++ libc-dev
WORKDIR /
 
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
 
RUN go mod download
 
COPY ./*.go /
RUN go build -o /expenses_bot
 
FROM alpine AS weaviate
COPY --from=builder /expenses_bot /expenses_bot
CMD ["./expenses_bot"]