# Security Guidelines for Vibe Agentic

This document contains security best practices and requirements for all developers working on the Vibe Agentic project.

## Input Validation and Sanitization

### Requirements

**ALL user input must be validated and sanitized** before processing. This includes:
- HTTP request bodies
- URL parameters
- Headers
- File uploads (when implemented)

### Implementation Rules

#### 1. JSON Input Validation

```go
// GOOD: Validate JSON structure and content
type LoginRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
    Password string `json:"password" validate:"required,min=8,max=100"`
}

var validator = validator.New()

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
}

if err := validator.Struct(req); err != nil {
    http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
    return
}
```

#### 2. Input Sanitization

```go
// GOOD: Sanitize input to prevent XSS
import "github.com/microcosm-cc/bluemonday"

func sanitizeInput(input string) string {
    p := bluemonday.UGCPolicy()
    return p.Sanitize(input)
}

// Use for all text that will be displayed
note.Text = sanitizeInput(requestBody.Text)
```

#### 3. Length Validation

```go
// GOOD: Prevent DoS attacks with length limits
const (
    MaxUsernameLength = 50
    MaxPasswordLength = 100
    MaxNoteTextLength = 10000
)

if len(req.Username) > MaxUsernameLength {
    http.Error(w, "Username too long", http.StatusBadRequest)
    return
}
```

#### 4. Content-Type Validation

```go
// GOOD: Validate Content-Type header
if r.Header.Get("Content-Type") != "application/json" {
    http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
    return
}
```

### Forbidden Practices

```go
// BAD: No validation at all
var req LoginRequest
json.NewDecoder(r.Body).Decode(&req) // ❌ UNSAFE

// BAD: Only checking decode success
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}
// ❌ Still UNSAFE - no content validation

// BAD: No sanitization before storage/display
note.Text = requestBody.Text // ❌ XSS risk
```

## Authentication Security

### Password Requirements

- Minimum 12 characters

```go
func isStrongPassword(password string) bool {
    return len(password) > 12
}
```

### Session Management

- Session IDs must be 32+ bytes of cryptographic random data
- Maximum session lifetime: 24 hours
- Sessions must be invalidated on password change
- Use `crypto/rand` for all security-sensitive random values

```go
// GOOD: Secure session ID generation
func generateSessionID() string {
    randomBytes := make([]byte, 32)
    if _, err := rand.Read(randomBytes); err != nil {
        // Handle error appropriately
        return ""
    }
    return hex.EncodeToString(randomBytes)
}
```

## Security Headers

All HTTP responses must include these security headers:

```go
w.Header().Set("Content-Security-Policy", "default-src 'self'")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
```

## Rate Limiting

All authentication endpoints must have rate limiting:

```go
// Maximum 5 login attempts per minute per IP
func (rl *RateLimiter) AllowRequest(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    count := rl.attempts[ip]
    if count >= 5 {
        return false
    }
    
    rl.attempts[ip] = count + 1
    // Reset counter after 1 minute
    go func() {
        time.AfterFunc(1*time.Minute, func() {
            rl.mu.Lock()
            delete(rl.attempts, ip)
            rl.mu.Unlock()
        })
    }()
    
    return true
}
```

## Error Handling

Never expose internal errors to clients:

```go
// GOOD: Sanitize errors
func handleError(w http.ResponseWriter, err error) {
    log.Printf("ERROR: %v", err) // Log internally
    
    // Return generic error to client
    http.Error(w, "An error occurred", http.StatusInternalServerError)
}

// BAD: Exposing internal errors
http.Error(w, err.Error(), http.StatusInternalServerError) // ❌ UNSAFE
```

## Logging and Monitoring

### Security Events to Log

1. **Authentication attempts** (success/failure)
2. **Session creation/destruction**
3. **Sensitive operations** (password changes, user management)
4. **Failed input validation**
5. **Rate limit triggers**

```go
// Minimum security log structure
type SecurityLog struct {
    Timestamp   time.Time
    EventType   string  // "login_attempt", "session_created", etc.
    Username    string
    IPAddress   string
    Success     bool
    Details     map[string]interface{}
}
```

## Dependency Security

1. **Regular dependency scanning** with `go mod tidy`
2. **Monthly dependency updates**
3. **Vulnerability monitoring** using `govulncheck`
4. **No development dependencies in production**

## Code Review Checklist

- [ ] All user input is validated
- [ ] Input is sanitized before storage/display
- [ ] Length limits are enforced
- [ ] Security headers are set
- [ ] Errors are properly sanitized
- [ ] No sensitive data in logs
- [ ] Rate limiting is implemented on auth endpoints
- [ ] Password requirements are enforced
- [ ] Session management follows security best practices

## Libraries and Tools

### Recommended Security Libraries

```bash
# Input validation
go get github.com/go-playground/validator/v10

# HTML sanitization
go get github.com/microcosm-cc/bluemonday

# Security headers
go get github.com/gorilla/handlers

# Rate limiting
go get golang.org/x/time/rate
```

### Security Tools
2. **Dependency Scanning**: `go mod tidy` + `go mod verify`

## Compliance Requirements

All code must comply with:
- OWASP Top Ten 2021
- CWE/SANS Top 25
- Go Secure Coding Practices

## Reporting Security Issues

Security vulnerabilities should be reported immediately to the security team and must be fixed within 24 hours for critical issues.

## Training Requirements

All developers must complete:
- OWASP Top Ten training
- Secure coding in Go
- Annual security refresher

## Change Log

- **2024-05-14**: Initial document created focusing on input validation requirements