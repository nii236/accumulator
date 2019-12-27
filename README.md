# Accumulator

Project Accumulator will track attendance rates per minute per student for VRNihongo classes.

## Database

```bash
rm accumulator.db
./db-prepare.sh
go run main.go -db-seed
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
