# Dynasty Project - Claude Code Instructions

## Project Overview

**Dynasty** is a "neomonolith" Go backend application for residential complex management, facilitating communication between residents and security services. It demonstrates microservices-inspired architecture within a monolithic Go application.

**Key Features:**
- User management with family member support
- JWT-based authentication and session management
- Request management (guest access, taxi, delivery, cargo)
- Image upload and CDN integration
- Password recovery system
- Role-based access control (admin, service, guard, neighbor)
- Multi-language support (English, Russian, Ukrainian)

**Technology Stack:**
- **Go 1.26.3** with Chi router
- **PostgreSQL** via GORM
- **JWT** authentication
- **S3/DigitalOcean Spaces** for image storage
- **Docker** with multi-stage builds
- **Traefik** reverse proxy
- **GitHub Actions** CI/CD

---

## Architecture Patterns

### Layered "Neomonolith" Structure

Each feature is organized as a self-contained handler with three clear layers:

```
server/handlers/{feature}/
├── service.go          # Business logic, interfaces
├── entities.go         # Domain models
├── transport/
│   ├── http.go        # HTTP handlers
│   ├── dto.go         # Request/Response DTOs
│   └── http_test.go   # Endpoint tests
├── repo/
│   └── repo.go        # Database operations
└── service_test.go    # Service tests with mocks
```

**Dependency Flow:** Transport → Service → Repository

### Key Architectural Principles

1. **Interface-based Design** - Services depend on interfaces, not concrete implementations
2. **Dependency Injection** - All dependencies passed via constructors (no global state)
3. **Clean Separation** - HTTP concerns in transport/, business logic in service, data access in repo/
4. **Testability** - Interfaces allow easy mocking with `moq`

### Example Service Initialization Pattern

```go
func New(log logger.Logger, repo userRepository, mail mailSender) *Service {
    return &Service{log: log, repo: repo, mail: mail}
}
```

---

## Critical Code Conventions

### Package Organization

- **Feature-based packages**: `auth`, `users`, `requests`, `dictionaries`, `health`
- **Sub-packages**: Always use `transport/` and `repo/` for separation
- **Lowercase package names**: Follow Go conventions

### Naming Standards

**Interfaces:**
- Lowercase with descriptive names: `userRepository`, `mailSender`, `s3Client`
- Defined locally in packages that use them (not exported)

**Structs:**
- Exported types: PascalCase (`User`, `Request`, `Session`)
- DTOs: Suffixed with purpose (`LoginRequest`, `UserByIDResponse`)
- Internal types: camelCase (`userRegisterRequest`, `errorResponse`)

**Methods:**
- HTTP handlers: Action-based (`Login()`, `Register()`, `Update()`)
- Service methods: Descriptive with context (`UserByPhoneAndPassword()`, `RegisterFamilyMember()`)
- Repository methods: CRUD operations (`CreateUser()`, `GetUserByID()`)

**Files:**
- Tests: `*_test.go`
- Mocks: `mock_test.go` (moq-generated)
- Service files: `service.go`, `service_master.go` (for extended logic)
- Transport: `http.go`, `dto.go`, `http_test.go`

### Code Quality Standards

- **Error Handling**: Always return errors explicitly, never panic
- **Context Propagation**: Pass context through call chain
- **Pointer Receivers**: Use for mutating methods
- **Linter Directives**: Use `nolint` for complex but necessary functions (`nolint: gocyclo,funlen`)
- **Security Directives**: Use `#nosec` for gosec false positives with explanatory comments (`#nosec G120 -- explanation`)

---

## Testing Practices

### Test Structure

- **30 test files** with **93 test functions**
- **Coverage tracking** via codecov.io
- **Race detection** enabled in Docker tests (`-race` flag)

### Testing Framework

```go
import (
    "github.com/stretchr/testify/assert"
    // Mocks generated with: go generate ./...
)
```

### Mock Generation

```bash
make gen  # Generates mocks using matryer/moq
```

**Pattern:**
- Service tests mock repositories
- Transport tests mock services
- Mocks defined in `mock_test.go` files

### Running Tests

```bash
make test        # Run tests with coverage
make cover       # Generate HTML coverage report
```

---

## Error Handling Pattern

### Centralized Error Definitions

All errors defined in `common/errs/errs.go`:

```go
type SvcError struct {
    Code int       // Error code (100+)
    Err  error     // English message
    Ru   string    // Russian translation
    Ua   string    // Ukrainian translation
}
```

### HTTP Error Response Format

```json
{
  "error": "English message",
  "error_code": 101,
  "ru": "Russian message",
  "ua": "Ukrainian message"
}
```

### Creating New Errors

When adding new error types:
1. Define in `common/errs/errs.go` with unique code
2. Provide all three language translations
3. Follow existing naming pattern: `{feature}{Reason}Code`

---

## Database Models

### Key Tables

- **users** - User accounts with roles, family hierarchy (parent_id)
- **user_roles** - Role definitions with hierarchy
- **buildings** & **entries** - Physical structure
- **sessions** - JWT refresh tokens
- **requests** - Service requests with images and status history
- **password_recovery** - Password reset tokens

### Entity Patterns

- **GORM Tags**: Use for column mapping and constraints
- **JSON Tags**: Always provide for API serialization
- **Soft Deletes**: Use `deleted_at` for safe deletion
- **Relationships**: Define in GORM structs (BelongsTo, HasMany)

### Schema Location

Database schema: `schema.sql`

---

## Configuration Management

### Environment Variables

Template: `cmd/.env.dist`

**Critical Variables:**
- `DB_*` - Database connection
- `S3_*` - Object storage (DigitalOcean Spaces)
- `SMTP_*` - Email service
- `AUTH_JWT_SECRET` - JWT signing key
- `LOG_LEVEL` - Logging verbosity

### Config Structure

Defined in `config/config.go`:
- Embedded structs for feature groups
- Validation on load
- Sensible defaults

---

## Security Best Practices

### Current Security Measures

1. **JWT Authentication** - Token-based with refresh mechanism
2. **Password Hashing** - bcrypt with 10 rounds
3. **HTML Sanitization** - bluemonday policy prevents XSS
4. **Input Validation** - go-playground/validator with struct tags
5. **Role-based Access** - User roles control endpoint access
6. **Session Management** - UUID-based refresh tokens with expiration
7. **Password Recovery** - Time-limited codes (3 hours)
8. **Bounded File Uploads** - MaxBytesReader + bounded ParseMultipartForm (10MB limit)
9. **Path Traversal Protection** - Multi-layered validation for static asset serving

### File Upload Security

**Location**: `server/handlers/requests/transport/http.go`

```go
const maxUploadSize = 10 << 20 // 10 MB

// Protect against unbounded uploads
r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+512)
// #nosec G120 -- ParseMultipartForm is bounded by maxUploadSize
if err := r.ParseMultipartForm(maxUploadSize); err != nil {
    // handle error
}
```

**Key Points:**
- Always use `http.MaxBytesReader` before parsing multipart forms
- Define upload size limits as constants
- Validate file size after parsing: `header.Size > (5 << 20)`
- Use `#nosec` directives with explanations for false positives

### Path Traversal Protection

**Location**: `server/handlers/ui/transport/http.go`

Multi-layered defense for static asset serving:

```go
// 1. Clean the path
cleanFilename := filepath.Clean(filename)

// 2. Reject absolute paths and traversal sequences
if filepath.IsAbs(cleanFilename) || strings.Contains(cleanFilename, "..") {
    w.WriteHeader(http.StatusBadRequest)
    return
}

// 3. Build safe path
basePath := filepath.Join("..", "_ui", "guard")
fullPath := filepath.Join(basePath, folder, cleanFilename)

// 4. Verify resolved path is within base directory
absBase, _ := filepath.Abs(basePath)
absFullPath, _ := filepath.Abs(fullPath)
if !strings.HasPrefix(absFullPath, absBase) {
    w.WriteHeader(http.StatusBadRequest)
    return
}
```

**Defense Layers:**
1. Path normalization with `filepath.Clean()`
2. Explicit rejection of absolute paths
3. Explicit rejection of paths containing ".."
4. Absolute path verification using `strings.HasPrefix()`

### Security Patterns to Follow

- **Never log passwords** or sensitive tokens
- **Validate all user input** using validator tags
- **Sanitize HTML** before storing user-generated content
- **Use prepared statements** (GORM handles this)
- **Check authorization** before resource access
- **Hash passwords** immediately on input
- **Bound all file operations** with size limits
- **Sanitize file paths** with multi-layered validation
- **Use #nosec directives** with clear explanations when security tools produce false positives
- **Always close file handles** with `defer fp.Close()`

---

## Common Development Tasks

### Building and Testing

```bash
make deps         # Download dependencies
make test         # Run tests with coverage
make lint         # Run golangci-lint
make build        # Build Docker image
make gen          # Generate mocks
```

### Docker Workflow

```bash
# Development database
docker-compose -f docker-database.yml up

# Production build and run
docker-compose up --build

# View logs
docker-compose logs -f backend
```

### Adding New Features

**Steps:**
1. Create handler directory: `server/handlers/{feature}/`
2. Define entities in `entities.go`
3. Create repository interface and implementation in `repo/`
4. Implement service with business logic in `service.go`
5. Add HTTP transport in `transport/http.go` and `transport/dto.go`
6. Write tests: `service_test.go`, `transport/http_test.go`
7. Generate mocks: `make gen`
8. Wire service in `cmd/main.go`
9. Mount routes in `server/http_server.go`

### Middleware Implementation

Implement the `Middleware` interface:

```go
type Middleware interface {
    Middleware(next http.Handler) http.Handler
}
```

Add to `server/middlewares/` and apply in `server/http_server.go`.

---

## Important Files and Locations

### Entry Points

- `cmd/main.go` - Application initialization and dependency wiring
- `server/http_server.go` - Router setup and middleware configuration

### Configuration

- `config/config.go` - Configuration struct and loading
- `cmd/.env.dist` - Environment variable template

### Core Logic

- `server/handlers/*` - Feature implementations
- `common/errs/errs.go` - Error definitions

### Infrastructure

- `Dockerfile` - Multi-stage build with tests
- `docker-compose.yml` - Production stack
- `Makefile` - Build automation
- `.github/workflows/` - CI/CD pipelines
- `_traefik/` - Reverse proxy configuration

### Database

- `schema.sql` - Database schema definition
- `server/handlers/*/repo/` - Database access layer

---

## Code Review Checklist

When reviewing or writing code, ensure:

- [ ] Follows layered architecture (transport/service/repo)
- [ ] Interfaces defined for dependencies
- [ ] Errors returned explicitly (no panics)
- [ ] Tests written with appropriate mocks
- [ ] Input validation using struct tags
- [ ] Error responses include multilingual messages
- [ ] No sensitive data in logs
- [ ] GORM used for all database operations
- [ ] Context propagated through call chain
- [ ] Code passes `make lint`
- [ ] File uploads are bounded with MaxBytesReader
- [ ] File paths are sanitized to prevent path traversal
- [ ] File handles closed with defer
- [ ] #nosec directives include explanatory comments

---

## CI/CD Pipeline

### GitHub Actions Workflows

**Master Branch** (`.github/workflows/main.yml`):
1. Run tests with coverage
2. Run golangci-lint
3. Build Docker image
4. Push to Docker Hub (ivch/dynasty)

**Pull Requests** (`.github/workflows/pr.yml`):
1. Validate tests pass
2. Check linting
3. Verify build succeeds

### Code Coverage

- Uploaded to codecov.io automatically
- View reports at https://codecov.io/gh/ivch/dynasty

---

## API Documentation

- **Postman Collection**: Maintained separately (mentioned in README)
- **API Versioning**: Prefix routes with `/v1/`
- **Authentication**: Bearer token in Authorization header

### Standard Response Formats

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

---

## Working with This Codebase

### Key Design Decisions

1. **Monolith over Microservices** - Simpler deployment, easier development
2. **Interface-based** - Enables future service extraction if needed
3. **GORM over Raw SQL** - Balance between abstraction and control
4. **Chi Router** - Lightweight, composable, idiomatic
5. **JWT + Refresh Tokens** - Stateless auth with long-lived sessions

### Future Extensibility

The architecture supports:
- Feature extraction to separate services
- API versioning (v2, v3, etc.)
- Multiple authentication methods
- Additional databases (GORM supports many)
- Horizontal scaling behind Traefik

### Performance Considerations

- **Database indexes** - Defined in schema.sql
- **Image processing** - Async upload to S3
- **Rate limiting** - Configured in Traefik
- **Connection pooling** - Managed by GORM

---

## Useful Commands Reference

```bash
# Development
make test                    # Run tests
make lint                    # Lint code
make cover                   # Generate coverage report
make gen                     # Generate mocks
go run cmd/main.go          # Run locally

# Docker
docker-compose up            # Start services
docker-compose logs -f       # View logs
docker-compose down          # Stop services

# Database
docker-compose -f docker-database.yml up  # Start dev DB

# Git workflow
git checkout -b feature/...  # Create feature branch
# Make changes, test, commit
# CI runs on PR creation
```

---

## Contact and Resources

- **GitHub**: https://github.com/ivch/dynasty
- **Docker Hub**: ivch/dynasty
- **Codecov**: https://codecov.io/gh/ivch/dynasty

---

## Notes for Claude Code

When working with this codebase:

1. **Always follow the layered architecture** - Don't mix transport/service/repo concerns
2. **Use interfaces for dependencies** - Never import concrete types from other handlers
3. **Maintain error translations** - All three languages (EN/RU/UA) required
4. **Test everything** - Write service and transport tests for new features
5. **Run make lint** - Before committing changes
6. **Generate mocks** - After changing interfaces (`make gen`)
7. **Update schema.sql** - When adding database tables/columns
8. **Follow existing patterns** - Consistency is key in this codebase
9. **Security first** - Validate input, sanitize output, hash passwords, bound file operations, sanitize paths
10. **Document breaking changes** - API changes affect clients
11. **File handling** - Always use MaxBytesReader, close file handles, validate paths
12. **Gosec compliance** - Add #nosec directives with explanations for false positives

### Security Checklist for New Code

When adding file handling:
- [ ] Use `http.MaxBytesReader` before parsing multipart forms
- [ ] Define size limits as constants
- [ ] Validate file size after parsing
- [ ] Close file handles with `defer fp.Close()`

When serving static files or handling user-provided paths:
- [ ] Clean paths with `filepath.Clean()`
- [ ] Reject absolute paths with `filepath.IsAbs()`
- [ ] Reject paths containing `..`
- [ ] Verify resolved path is within base directory

When adding features, reference existing handlers (e.g., `users`, `requests`) as templates for structure and patterns.