# Build stage
FROM golang:1.24 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# build dari folder ./cmd
RUN go build -o server ./cmd

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build /app/server .

CMD ["/app/server"]
