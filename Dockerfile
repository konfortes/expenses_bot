FROM golang:1.11-alpine AS build_base
 
RUN apk add bash git gcc g++ libc-dev
WORKDIR /
 
ENV GO111MODULE=on
 COPY go.mod .
COPY go.sum .
 
RUN go mod download
 
FROM build_base AS server_builder
COPY . .
RUN go build -o /expenses_go
 
FROM alpine AS weaviate
COPY --from=server_builder /expenses_go /expenses_go
ENTRYPOINT ["/expenses_go"]