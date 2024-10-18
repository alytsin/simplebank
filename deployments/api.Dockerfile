FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/server.go


FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY config.env .
EXPOSE 8080
CMD [ "/app/server" ]