# Accumulator

Project Accumulator will track attendance rates per minute per student for VRNihongo classes.

## Database

```bash
rm accumulator.db
# Embed migration bindata (ignore errors)
go generate
go run cmd/admin/main.go -db-migrate
# Generate SQLboiler
go generate
go run cmd/accumulator/main.go -db-seed
```

## Server

```bash
go run main.go
```

## Frontend

```bash
cd web
npm install
npm start
```
