# Gunakan image golang resmi
FROM golang:1.21 as builder

WORKDIR /app
COPY . .

RUN go build -o server .

# Stage kedua: image ringan
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
