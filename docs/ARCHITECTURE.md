# Architecture Overview

## General
This is a Go language project.

## Control
The project consists of a REST API and a CLI for administrative actions.

## Module Structure
The system is organized into subdomains. Each subdomain module follows these rules:

### Layers within each Module:
1. **Domain Layer**: Contains business objects and core logic
2. **Storage Layer**: Handles data persistence
3. **Service Layer**: Contains complex business logic
4. **Controller Layer**: REST handlers and CLI commands

### Ports and Adapters Pattern
Each module follows a ports and adapters architecture:
- Domain layer defines interfaces (ports)
- Infrastructure implements interfaces (adapters)
- Dependencies point inward toward the domain

### File Organization per Module:
- `rest.go` - REST controllers
- `domain.go` - Domain objects with business logic
- `repository.go` - Repository interfaces
- `<type>_repository.go` - Repository implementations
- `service.go` - Complex business logic services
- `cmd/cli/cmd/<command>.go` - CLI commands

## Authorization ⚠️ IMPORTANT
Every API request must be authorized with a Session-ID. Authorization must meet these requirements:

1. **Session-ID Required**: Every request must contain a valid Session-ID in the header
2. **Endpoint-specific Permissions**: Each endpoint requires specific HTTP permissions
3. **Middleware-based**: Authorization must be implemented as middleware executed before handlers
4. **Error Handling**: Unauthorized requests must return 401 Unauthorized
5. **Session Storage**: Sessions must be stored and validated in a Session-Store

### Implementation Guidelines:
- Middleware file: `middleware/auth.go`
- Session management: `auth/session.go`
- Session-Store interface: `auth.SessionStore`
- In-Memory implementation: `auth.InMemorySessionStore`
- Permission check must be called before each handler
- Session-Store must be configured before handler registration

### Best Practices:
- ✅ Inject Session-Store as dependency
- ✅ Use thread-safe implementation (Mutex)
- ✅ Implement session timeout (e.g., 24 hours)
- ✅ Generate session IDs with UUID or similar
- ✅ Validate sessions against the store

### Anti-Patterns:
- ❌ Manual session management
- ❌ Validate session IDs without store
- ❌ Use Session-Store globally without dependency injection
- ❌ Ignore thread safety

### Example:
```go
// Correct implementation
sessionStore := auth.NewInMemorySessionStore()
auth.RegisterHandlers(sessionStore)
middleware.SetSessionStore(sessionStore)

// Anti-Pattern (wrong implementation)
// auth.RegisterHandlers() // Without Session-Store
```

## Domain Layer
The Domain layer contains business objects and their core logic. Each domain object should encapsulate its own business logic.

### User Domain Object
The User object in `auth/user/domain.go` contains:
- Password hashing (`HashPassword()`)
- Password verification (`VerifyPassword()`)
- Password change (`SetPassword()`)

### Best Practices for Domain Objects:
- ✅ Domain objects should contain their own validation logic
- ✅ Business rules should be encapsulated in domain methods
- ✅ Domain objects should be independent of storage or UI concerns
- ✅ Use constructors for complex object creation

### Anti-Patterns for Domain Objects:
- ❌ Move domain logic to repositories or services
- ❌ Create domain objects with storage dependencies
- ❌ Implement business rules in UI or CLI layer

## Storage Layer
Storage uses individual JSON files. Each file has a generic ID and is located in its entity's directory.
Example: `lists/alsdkjsfhakdsnjva.json`

## CLI Structure
Administrative functions are implemented as CLI commands:
- Command files follow the pattern: `cmd/cli/cmd/<command>.go`
- Each command has clear help documentation
- Commands return structured output (JSON or text)
- Errors are handled with appropriate exit codes and messages

### Best Practices for CLI Commands:
- ✅ Use Cobra or similar framework for CLI structure
- ✅ Implement help flags (--help)
- ✅ Provide clear success and error messages
- ✅ Use appropriate exit codes (0 for success, !=0 for errors)
- ✅ Support both JSON and text output
- ✅ Use domain objects for business logic (e.g., User.HashPassword())

### Anti-Patterns for CLI:
- ❌ Direct business logic in CLI files
- ❌ No error handling in commands
- ❌ Unclear or missing documentation
- ❌ Inconsistent output formats
- ❌ Manual implementation of domain logic (e.g., password hashing in CLI)

## Checklist for New Endpoints
1. [ ] Implement/integrate authorization middleware
2. [ ] Implement session-ID validation
3. [ ] Check endpoint-specific permissions
4. [ ] Handle error cases (401, 403) correctly