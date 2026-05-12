package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session represents a user session
type Session struct {
	ID        string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore interface for session management
type SessionStore interface {
	CreateSession(username string) (Session, error)
	GetSession(sessionID string) (Session, bool)
	DeleteSession(sessionID string) bool
}

// InMemorySessionStore implements SessionStore in memory
type InMemorySessionStore struct {
	sessions map[string]Session
	mu       sync.RWMutex
}

// NewInMemorySessionStore creates a new in-memory session store
func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]Session),
	}
}

// CreateSession creates a new session for the given username
func (s *InMemorySessionStore) CreateSession(username string) (Session, error) {
	session := Session{
		ID:        generateSessionID(),
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiration
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	
	return session, nil
}

// GetSession retrieves a session by ID
func (s *InMemorySessionStore) GetSession(sessionID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	session, exists := s.sessions[sessionID]
	if !exists {
		return Session{}, false
	}
	
	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return Session{}, false
	}
	
	return session, true
}

// DeleteSession removes a session by ID
func (s *InMemorySessionStore) DeleteSession(sessionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.sessions[sessionID]; exists {
		delete(s.sessions, sessionID)
		return true
	}
	return false
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	// Generate a random 32-byte session ID (64 characters in hex)
	// Using crypto/rand for cryptographic security
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to timestamp if crypto/rand fails
		return time.Now().Format("20060102150405.999999999")
	}
	return hex.EncodeToString(randomBytes)
}
