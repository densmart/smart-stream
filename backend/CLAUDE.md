# smart-stream Backend (Go)

This is the backend service for the smart-stream video streaming platform, written in Go.

## Architecture Overview

Clean architecture with separation of concerns:

```
backend/
├── cmd/                    # Application entry points
│   └── main.go            # Main server application
├── internal/
│   ├── adapters/          # External interface adapters
│   │   ├── api/          # HTTP handlers (Gin)
│   │   ├── db/           # Database implementations
│   │   └── dto/          # Data Transfer Objects
│   └── domain/           # Core business logic
│       ├── models/       # Domain models
│       ├── repo/         # Repository interfaces
│       ├── usecases/     # Business use cases
│       └── utils/        # Domain utilities
├── config/               # Configuration files
├── migrations/           # Database migrations (tern)
└── builds/              # Build artifacts
```

## Key Technologies

- **Web Framework**: Gin
- **Database**: PostgreSQL (jackc/pgx/v4)
- **Migrations**: jackc/tern
- **Authentication**: JWT (separate TTL for web/TV clients)
- **Config**: Viper
- **Logging**: Uber Zap
- **Query Builder**: goqu
- **Video Processing**: go-ffprobe

## Database

- PostgreSQL as primary data store
- Migrations managed with tern
- Connection pooling via pgx
- Transactions for data consistency

## API Structure

- RESTful endpoints via Gin
- JWT authentication middleware
- Request validation
- Structured error responses
- CORS configuration

## Configuration

Environment variables managed through:
- `.env` file (local development)
- Viper for configuration loading
- Separate config for different environments

## Running Locally

```bash
cd backend
go run cmd/main.go
```

## Database Migrations

```bash
# Run migrations
tern migrate -m migrations/

# Create new migration
tern new -m migrations/ migration_name
```

## Testing

```bash
go test ./...
```

## Building

```bash
# Development build
go build -o builds/server cmd/main.go

# Production build (see Dockerfile)
```

## API Authentication

- JWT tokens for user authentication
- Separate token TTLs:
  - Web client: Standard TTL
  - TV client: Extended TTL
- Token validation middleware
- Secure token generation

## Video Streaming

- Video metadata extraction using ffprobe
- Streaming optimization
- Format compatibility checks
- Thumbnail generation

## Important Notes

- Follow clean architecture principles
- Keep domain logic independent of infrastructure
- Use repository pattern for all data access
- Write migrations for any schema changes
- Log important operations with structured logging
- Validate all user inputs
- Handle errors gracefully with meaningful messages
