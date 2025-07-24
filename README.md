# Parking Service

A modular Go backend that can run standalone or be imported into other applications. Includes:
- REST API for managing parking spots
- Google Maps geocoding
- HTML/JS frontend map viewer

## Run Locally
```bash
git clone <repo>
cd github.com/khanjumed/Geopackage
cp .env.example .env
go run cmd/server/main.go
```

## Use in Another Go App
```go
import "github.com/yourusername/github.com/khanjumed/Geopackage/internal/routes"
r := gin.Default()
routes.Register(r)
r.Run()
```