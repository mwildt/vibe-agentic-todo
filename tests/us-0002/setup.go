package tests

import (
	"sync"
	"vibe-agentic/auth"
	"vibe-agentic/middleware"
)

var (
	setupOnce sync.Once
	sessionStore = auth.NewInMemorySessionStore()
)

func setupTest() {
	setupOnce.Do(func() {
		// Register auth handlers with session store
		auth.RegisterHandlers(sessionStore)
		
		// Configure middleware with session store
		middleware.SetSessionStore(sessionStore)
	})
}
