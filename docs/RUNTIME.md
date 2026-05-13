# Runtime Information

## Session Management

### Session Lifecycle
- Sessions are created upon successful login via `/login` endpoint
- Session IDs are 64-character hex-encoded 32-byte random values
- Sessions are stored in an in-memory session store
- Default session timeout: 24 hours

### Session Storage
- Implementation: `auth.InMemorySessionStore`
- Thread-safe using `sync.Mutex`
- Stores mapping of sessionID → username

### Session Validation
1. Middleware extracts `session_id` cookie
2. Validates session exists in store
3. Rejects with 401 Unauthorized if invalid/missing
4. For valid sessions, adds username to request context

## Authentication Flow

```
Client → POST /login (username, password) → Server
       ← 200 OK (session_id) ←

Client → GET /protected (X-Session-ID: session_id) → Server
       ← 200 OK (protected data) ←
```

### Login Process
1. Client sends POST request to `/login` with JSON body:
   ```json
   {
     "username": "user",
     "password": "password"
   }
   ```
2. Server validates credentials against user repository
3. On success: generates session ID, stores in session store, sets HTTP cookie, returns 200
4. On failure: returns 401 Unauthorized

### Protected Requests
1. Client includes `session_id` cookie with valid session ID
2. Auth middleware validates session
3. On success: adds username to context, calls handler
4. On failure: returns 401 Unauthorized

## Error Handling

### HTTP Status Codes
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input/validation failed
- `401 Unauthorized`: Missing/invalid session
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

### Error Response Format
```json
{
  "error": "Error message",
  "details": "Additional details"
}
```

## Security Considerations

### Password Storage
- Passwords are hashed using bcrypt
- Default cost factor: bcrypt.DefaultCost
- Hashes are stored in user repository

### Session Security
- Session IDs are cryptographically secure random values
- 32-byte (256-bit) random values, hex-encoded to 64 characters
- Session timeout prevents indefinite access
- No session fixation vulnerabilities
- HttpOnly cookies prevent JavaScript access
- Secure flag for HTTPS-only cookies in production
- SameSite attributes prevent CSRF attacks

### Input Validation
- All API inputs are validated
- Required fields are enforced
- Data types are validated
- String lengths are checked

## Performance Characteristics

### Session Operations
- Session creation: O(1)
- Session validation: O(1)
- Session storage: In-memory hash map

### API Response Times
- Login: < 100ms (with bcrypt hashing)
- Protected endpoints: < 50ms (with auth middleware)
- Note operations: < 30ms (JSON file I/O)

### Concurrency
- Thread-safe session store using mutex
- Concurrent request handling supported
- No blocking operations in critical path

## Configuration

### Environment Variables
- `VIBe_HOME`: Base directory for data storage
- `SESSION_TIMEOUT`: Session timeout duration (default: 24h)
- `BCRYPT_COST`: Bcrypt cost factor (default: bcrypt.DefaultCost)

### Data Directories
- User data: `$VIBe_HOME/users.yaml`
- Note data: `$VIBe_HOME/notes/`
- Session data: In-memory (not persisted)

## Monitoring and Logging

### Security Events
The system logs security-related events:
- Login attempts (success/failure)
- Session creation
- Authentication failures

### Event Format
```
SECUITY_EVENT: type=event_type username=user ip= success=bool details=map[key:value]
```

### Example Events
```
SECURITY_EVENT: type=login_attempt username=testuser ip= success=true details=map[reason:successful_authentication]
SECURITY_EVENT: type=session_creation username=testuser ip= success=true details=map[session_id:abc123...]
SECURITY_EVENT: type=login_attempt username=wronguser ip= success=false details=map[reason:invalid_credentials]
```

## Deployment Considerations

### Scaling
- Session store is in-memory (single instance only)
- For horizontal scaling: replace with distributed session store
- Stateless design allows multiple instances behind load balancer

### Persistence
- User data persisted to YAML file
- Note data persisted to JSON files
- Sessions not persisted (lost on restart)

### Backup Strategy
- Regular backups of `$VIBe_HOME` directory
- User data: critical (contains password hashes)
- Note data: important (user content)

## Development Runtime

### Development Mode
- Run with `go run main.go`
- Automatically reloads on code changes (with appropriate tools)
- Debug logging enabled

### Production Mode
- Build with `go build -o vibe-server main.go`
- Run as service/user with limited privileges
- Enable proper logging and monitoring

### Testing Runtime
- Test data stored in `./test_data/`
- Automatically cleaned up after tests
- Isolated from production data
- Mock session store for testing