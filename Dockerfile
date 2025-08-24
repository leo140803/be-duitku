FROM golang:1.24 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server 


# runtime image lebih kecil
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/server .
CMD ["/app/server"]
