# Gunakan base image Golang
FROM golang:1.23 AS builder
WORKDIR /app

# Copy module dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy seluruh kode proyek
COPY . .

# Build aplikasi
RUN go build -o main .

# Gunakan base image yang lebih ringan
FROM gcr.io/distroless/base
WORKDIR /root/
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["/root/main"]
