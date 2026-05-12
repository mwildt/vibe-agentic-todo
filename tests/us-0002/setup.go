package tests

import (
	"sync"
	"vibe-agentic/auth"
	"vibe-agentic/auth/user"
	"vibe-agentic/middleware"
)

var (
	setupOnce sync.Once
	sessionStore = auth.NewInMemorySessionStore()
	userRepo      = user.NewYAMLUserRepository("users.yaml")
)

func setupTest() {
	setupOnce.Do(func() {
		// Load users
		userRepo.LoadUsers()
		
		// Create testuser if it doesn't exist
		if _, exists := userRepo.GetUser("testuser"); !exists {
			userRepo.CreateUser("testuser", "testpass")
		}
		
		// Register auth handlers with session store and user repository
		auth.RegisterHandlers(sessionStore, userRepo)
		
		// Configure middleware with session store
		middleware.SetSessionStore(sessionStore)
	})
}
