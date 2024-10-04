FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

COPY . .

RUN chmod +x /app/entrypoint.sh

RUN go build -o main .

# Menjalankan migrasi saat build atau container dijalankan
CMD ["/app/entrypoint.sh"]
