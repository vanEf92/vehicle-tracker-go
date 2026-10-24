Running at Windows 11 OS
Go version 1.27
Docker version 29.9.2

docker compose up --build -d
go run publisher/main.go
go run migrations/db_migrations.go

mqtt = localhost:1883

rabbitMQ = http://localhost:15672

docker compose down -v