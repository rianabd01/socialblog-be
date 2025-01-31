# Gunakan base image Go
FROM golang:1.20 AS builder

# Set working directory
WORKDIR /app

# Copy module dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy seluruh kode proyek
COPY . .

# Build aplikasi
RUN go build -o main .

# Gunakan base image yang lebih ringan untuk hasil akhir
FROM gcr.io/distroless/base-debian11
WORKDIR /root/
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["/root/main"]
