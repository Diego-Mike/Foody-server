# build serve
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o foody-server ./cmd/server

# run server
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/foody-server .
COPY .env .

EXPOSE 8000
CMD [ "/app/foody-server" ]