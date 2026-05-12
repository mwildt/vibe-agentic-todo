# Vibe Agentic - TODO List Management System

A modern TODO list management system with REST API and CLI administration.

## Features

### REST API
- **Session-based authentication** with JWT-like session tokens
- **TODO list management** (create, retrieve, update, delete)
- **User management** (create, update, delete)
- **Middleware-based authorization** for all endpoints

### CLI Administration
- **User management commands** (`user add`, `user update`, `user delete`)
- **Help system** with comprehensive documentation
- **JSON and text output** support

### Architecture Highlights
- **Domain-Driven Design** with clear separation of concerns
- **Clean Architecture** principles
- **Dependency Injection** throughout
- **Thread-safe** implementations
- **Comprehensive testing** coverage

## Getting Started

### Prerequisites
- Go 1.20+
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/vibe-agentic.git
cd vibe-agentic

# Install dependencies
go mod download
```

### Running the Server

```bash
# Start the REST API server
go run main.go

# The server will start on http://localhost:8080
```

### Using the CLI

```bash
# Build the CLI
go build -o vibe-cli ./cmd/cli

# Add a new user
./vibe-cli user add --username admin --password securepassword

# Update a user
./vibe-cli user update --username admin --password newpassword

# Delete a user
./vibe-cli user delete --username testuser

# Show help
./vibe-cli --help
```

## API Documentation

### Authentication

**Login Endpoint**
```
POST /login
Content-Type: application/json

{
  "username": "your_username",
  "password": "your_password"
}

Response:
{
  "session_id": "64-character-hex-string"
}
```

**Authenticated Requests**
All API endpoints require the `X-Session-ID` header:
```
GET /notes
X-Session-ID: your-session-id-here
```

### TODO Endpoints

**Create Note**
```
POST /notes
X-Session-ID: your-session-id
Content-Type: application/json

{
  "text": "Buy groceries"
}

Response: 201 Created
{
  "id": "generated-note-id",
  "text": "Buy groceries"
}
```

**Get Note**
```
GET /notes/{id}
X-Session-ID: your-session-id

Response: 200 OK
{
  "id": "note-id",
  "text": "Note content"
}
```

## Project Structure

```
auth/
├── rest.go            # Authentication REST handlers
├── session.go        # Session management
└── user/
    ├── domain.go     # User domain object with business logic
    └── repository.go  # User repository implementation

cmd/cli/
├── main.go           # CLI entry point
└── cmd/
    └── user.go        # User management commands

middleware/
└── auth.go           # Authentication middleware

notes/
├── rest.go           # Notes REST handlers
├── service.go        # Notes business logic
└── repository.go     # Notes repository interface

tests/
└── us-*/             # User story test suites

docs/
└── ARCHITECTURE.md   # Architecture documentation
```

## Architecture Principles

### Domain Layer
- **User domain object** encapsulates password hashing and verification
- **Business rules** are contained within domain objects
- **Pure domain logic** independent of storage or UI concerns

### Storage Layer
- **YAML-based user storage** for simplicity
- **JSON-based note storage** for flexibility
- **Repository pattern** for data access abstraction

### API Layer
- **RESTful endpoints** following standard conventions
- **Session-based authentication** with 24-hour expiration
- **Middleware-based authorization** for all protected routes

### CLI Layer
- **Cobra framework** for command structure
- **Domain object utilization** for business logic
- **Consistent output formats** (JSON/text)

## Development

### Running Tests

```bash
# Run all tests
cd tests && go test ./...

# Run specific test suite
cd tests && go test ./us-0001 -v
```

### Adding New Features

1. **Create domain objects** in `domain.go` files
2. **Implement repositories** following the repository pattern
3. **Add REST handlers** with proper authentication middleware
4. **Write tests** following existing patterns
5. **Update documentation** in `docs/ARCHITECTURE.md`

## Security

- **Bcrypt password hashing** with default cost factor
- **Session timeout** (24 hours)
- **CSRF protection** via session tokens
- **Input validation** on all endpoints

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature`
5. Open a pull request

## License

MIT License - see LICENSE file for details.

## Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

**Vibe Agentic** - Modern TODO management with clean architecture and comprehensive testing.