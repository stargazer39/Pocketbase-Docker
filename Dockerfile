FROM golang:1.21 as builder
WORKDIR /app
COPY . .
ARG GONOPROXY=*
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server server
EXPOSE 8090
CMD ["./server", "serve", "--http", ":8090"]