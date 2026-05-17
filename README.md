![Master](https://github.com/ivch/dynasty/workflows/Master/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/ivch/dynasty/branch/master/graph/badge.svg)](https://codecov.io/gh/ivch/dynasty)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivch/dynasty)](https://goreportcard.com/report/github.com/ivch/dynasty)

# Dynasty

A "neomonolith" Go backend application for residential complex management, facilitating communication between residents and security services.

## Disclaimer

This is my pet project to help my neighbors in communication with the security service of the residential complex.

The secondary goal of this project is to create a codebase that can represent my coding skills and, I hope, will help me to avoid test tasks on the interview :)

Read more about the architectural approach and motivations in [INFO.md](https://github.com/ivch/dynasty/blob/master/INFO.md).

## Features

- **User Management** - Registration, authentication, family member support
- **JWT Authentication** - Token-based auth with refresh mechanism
- **Request System** - Guest access, taxi, delivery, and cargo requests
- **Image Handling** - Upload and CDN integration via S3-compatible storage
- **Password Recovery** - Secure password reset with email verification
- **Role-Based Access** - Admin, service, guard, and neighbor roles
- **Multi-Language Support** - English, Russian, and Ukrainian error messages
- **Health Checks** - Monitoring endpoints for service health

## Tech Stack

- **Go 1.22-1.23** - Primary language
- **Chi Router** - Lightweight HTTP router
- **GORM** - ORM for PostgreSQL
- **JWT** - Token-based authentication
- **PostgreSQL** - Primary database
- **Docker** - Containerization with multi-stage builds
- **Traefik** - Reverse proxy and load balancer
- **GitHub Actions** - CI/CD pipeline
- **DigitalOcean Spaces** - S3-compatible object storage

## Architecture

Dynasty follows a **"neomonolith"** architecture - a monolithic application organized with microservices-inspired patterns. Each feature is self-contained with clear separation of concerns:

```
server/handlers/{feature}/
├── service.go          # Business logic
├── entities.go         # Domain models
├── transport/          # HTTP handlers, DTOs
├── repo/              # Database operations
└── *_test.go          # Tests with mocks
```

**Key Principles:**
- Interface-based design with dependency injection
- Clean separation of transport, service, and repository layers
- Comprehensive test coverage with mock generation
- Production-ready deployment configuration

## Prerequisites

- **Go 1.22+**
- **Docker & Docker Compose**
- **PostgreSQL** (via Docker or local)
- **Make** (for build automation)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/ivch/dynasty.git
cd dynasty
```

### 2. Configure Environment

Copy the environment template and configure your settings:

```bash
cp cmd/.env.dist cmd/.env
```

Edit `cmd/.env` with your configuration:
- Database connection settings
- S3/CDN credentials (DigitalOcean Spaces or AWS S3)
- SMTP settings for email
- JWT secret
- Log level

### 3. Start Development Database

```bash
docker-compose -f docker-database.yml up -d
```

This starts PostgreSQL on port 5432.

### 4. Run Locally (Development)

```bash
# Install dependencies
make deps

# Run the application
go run cmd/main.go
```

The server will start on the port specified in your `.env` file (default: 9001).

### 5. Run with Docker (Production)

```bash
# Build and start all services
docker-compose up --build

# Or build separately
make build
docker-compose up
```

This starts:
- Traefik reverse proxy
- Dynasty backend service
- Automatic SSL/TLS with Let's Encrypt

## Development

### Project Structure

```
dynasty/
├── cmd/                        # Application entry point
├── server/                     # HTTP server and handlers
│   ├── handlers/              # Feature modules (auth, users, requests, etc.)
│   ├── middlewares/           # HTTP middlewares
│   └── http_server.go         # Router setup
├── config/                     # Configuration management
├── common/                     # Shared utilities (errors, logger, email)
├── _traefik/                   # Traefik configuration
├── _ui/                        # Frontend assets
├── schema.sql                  # Database schema
├── Dockerfile                  # Multi-stage build
├── docker-compose.yml          # Production stack
├── Makefile                    # Build automation
└── CLAUDE.md                   # Development guide for AI assistants
```

### Available Make Commands

```bash
make deps         # Download and vendor dependencies
make test         # Run tests with coverage
make lint         # Run golangci-lint
make cover        # Generate HTML coverage report
make build        # Build Docker image with tests
make gen          # Generate mocks using moq
make tag          # Tag Docker image
make push         # Push to Docker Hub
```

### Running Tests

```bash
# Run all tests with coverage
make test

# Generate coverage report
make cover

# Run specific package tests
go test ./server/handlers/users/...
```

Current coverage: [![codecov](https://codecov.io/gh/ivch/dynasty/branch/master/graph/badge.svg)](https://codecov.io/gh/ivch/dynasty)

### Code Quality

```bash
# Run linter
make lint

# Format code
go fmt ./...

# Vet code
go vet ./...
```

Linting configuration: `.golangci.yml`

### Adding New Features

1. Create handler directory: `server/handlers/{feature}/`
2. Define domain models in `entities.go`
3. Create repository interface and implementation in `repo/`
4. Implement business logic in `service.go`
5. Add HTTP transport in `transport/` with DTOs and handlers
6. Write tests: `service_test.go` and `transport/http_test.go`
7. Generate mocks: `make gen`
8. Wire service in `cmd/main.go`
9. Mount routes in `server/http_server.go`

See `CLAUDE.md` for detailed development guidelines.

## API Documentation

Full API documentation available on Postman:

**[View API Documentation](https://documenter.getpostman.com/view/712107/SzKVSeBL?version=latest)**

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### Standard Response Format

**Success:**
```json
{
  "data": { ... }
}
```

**Error:**
```json
{
  "error": "English message",
  "error_code": 101,
  "ru": "Russian message",
  "ua": "Ukrainian message"
}
```

## Database

### Schema

Database schema is defined in `schema.sql`. Key tables:

- **users** - User accounts with roles and family hierarchy
- **user_roles** - Role definitions (admin, service, guard, neighbor)
- **buildings & entries** - Physical structure mapping
- **sessions** - JWT refresh tokens
- **requests** - Service requests with images and history
- **password_recovery** - Password reset tokens

### Migrations

Currently using schema.sql for initialization. For production, consider adding migration tool like [golang-migrate](https://github.com/golang-migrate/migrate).

### Development Database

```bash
# Start PostgreSQL
docker-compose -f docker-database.yml up -d

# Access database
psql -h localhost -p 5432 -U postgres -d dynasty
```

## Deployment

### Docker Production Build

The Dockerfile uses multi-stage builds:

1. **Build stage** - Compiles Go binary, runs tests with race detection
2. **Test stage** - Uploads coverage to codecov.io
3. **Final stage** - Minimal `scratch` image with binary only

```bash
# Build image
docker build -t dynasty:latest .

# Run with docker-compose
docker-compose up -d
```

### Environment Variables

Required environment variables in production:

```
DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_SCHEMA
S3_KEY, S3_SECRET, S3_ENDPOINT, S3_SPACE_NAME, CDN_HOST
SMTP_HOST, SMTP_PORT, SMTP_FROM, SMTP_PASS
AUTH_JWT_SECRET
HTTP_PORT
LOG_LEVEL
```

See `cmd/.env.dist` for complete list.

### Traefik Configuration

Traefik handles:
- Reverse proxy routing
- SSL/TLS termination with Let's Encrypt
- Rate limiting
- Compression
- Load balancing

Configuration in `_traefik/` directory.

## CI/CD

GitHub Actions workflows:

- **Master Branch** - Build, test, lint, push to Docker Hub
- **Pull Requests** - Validate tests and linting

Workflows: `.github/workflows/`

## Contributing

This is primarily a personal project, but suggestions and improvements are welcome!

### Contribution Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow existing code patterns and architecture
4. Write tests for new features
5. Ensure `make test` and `make lint` pass
6. Commit changes with clear messages
7. Push to your fork and create a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Maintain layered architecture (transport/service/repo)
- Use interfaces for dependencies
- Write tests with mocks
- Add multilingual error messages (EN/RU/UA)
- Keep functions focused and testable

See `CLAUDE.md` for detailed coding guidelines.

## License

This is a personal project. Use at your own risk.

## Contact

- **GitHub**: [@ivch](https://github.com/ivch)
- **Docker Hub**: [ivch/dynasty](https://hub.docker.com/r/ivch/dynasty)

## Acknowledgments

- Inspired by Alan Shreve's [Neomonolith](https://inconshreveable.com/10-07-2015/the-neomonolith/) article
- Built with excellent open-source Go libraries
- Thanks to the Go community for best practices and patterns

---

**Note**: This is a demonstration project showcasing Go backend development practices. While functional, additional hardening may be needed for production use at scale.
