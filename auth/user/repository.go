package user

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type User struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

type UserRepository interface {
	CreateUser(username, password string) error
	UpdateUser(username, password string) error
	DeleteUser(username string) error
	GetUser(username string) (User, bool)
}

type YAMLUserRepository struct {
	filePath string
	users    []User
	mu       sync.Mutex
}

func NewYAMLUserRepository(filePath string) *YAMLUserRepository {
	return &YAMLUserRepository{
		filePath: filePath,
		users:    []User{},
	}
}

func (r *YAMLUserRepository) LoadUsers() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the file with empty users
			err := r.saveUsers()
			return err
		}
		return err
	}

	var config struct {
		Users []User `yaml:"users"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	r.users = config.Users
	return nil
}

func (r *YAMLUserRepository) saveUsers() error {
	config := struct {
		Users []User `yaml:"users"`
	}{
		Users: r.users,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *YAMLUserRepository) CreateUser(username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	for _, user := range r.users {
		if user.Username == username {
			return fmt.Errorf("user '%s' already exists", username)
		}
	}

	// Add new user
	r.users = append(r.users, User{
		Username:     username,
		PasswordHash: hashPassword(password),
	})

	return r.saveUsers()
}

func (r *YAMLUserRepository) UpdateUser(username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find and update the user
	found := false
	for i, user := range r.users {
		if user.Username == username {
			r.users[i].PasswordHash = hashPassword(password)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("user '%s' not found", username)
	}

	return r.saveUsers()
}

func (r *YAMLUserRepository) DeleteUser(username string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find and remove the user
	newUsers := []User{}
	found := false
	for _, user := range r.users {
		if user.Username == username {
			found = true
		} else {
			newUsers = append(newUsers, user)
		}
	}

	if !found {
		return fmt.Errorf("user '%s' not found", username)
	}

	r.users = newUsers
	return r.saveUsers()
}

func (r *YAMLUserRepository) GetUser(username string) (User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, true
		}
	}

	return User{}, false
}

func hashPassword(password string) string {
	// Simple hash function for testing
	// In production, use bcrypt or similar
	return fmt.Sprintf("hashed_%s", password)
}
