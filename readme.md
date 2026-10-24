Running at Windows 11 OS
Go version 1.27
Docker version 29.9.2

Jalankan Docker dengan perintah:
docker compose up --build -d

Jalan script DB Migration untuk create table awal:
go run migrations/db_migrations.go

Jalankan script publisher:
go run publisher/main.go

**note: untuk pengujian yg lebih dalam, bisa mengganti posisi latitude dan longitude di file main.go dalam folder publisher

Endpoint Vehicle Location History:
http://localhost:8080/vehicles/:vehicle_id/history?start=1761282188&end=1824354188

Endpoint Vehicle Latest Location:
http://localhost:8080/vehicles/:vehicle_id/location

Postman Collection:
https://eevan-ef-823817.postman.co/workspace/Evan-Efendi's-Workspace~9b3479a5-3ff7-43d1-85ce-46dfdad99b5d/collection/59176886-8b7cd33c-8911-4390-98da-a3c5e01dcd1b?action=share&source=copy-link&creator=59176886

Akses MQTT menggunakan MQTT Explorer di localhost:1883

Akses Dashboard rabbitMQ = http://localhost:15672