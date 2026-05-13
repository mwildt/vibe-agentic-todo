package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system with password management capabilities
type User struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

// NewUser creates a new User instance
func NewUser(username, password string) (*User, error) {
	user := &User{
		Username: username,
	}
	
	// Hash the password
	hashedPassword, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}
	
	user.PasswordHash = hashedPassword
	return user, nil
}

// HashPassword hashes a plain text password and returns the bcrypt hash
func (u *User) HashPassword(password string) (string, error) {
	// Validate password length requirement (minimum 12 characters)
	if len(password) < 12 {
		return "", errors.New("password must be at least 12 characters long")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against the stored hash
func (u *User) VerifyPassword(password string) bool {
	if u.PasswordHash == "" {
		return false
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// SetPassword updates the user's password with proper hashing
func (u *User) SetPassword(password string) error {
	// Validate password length requirement (minimum 12 characters)
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	hashedPassword, err := u.HashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = hashedPassword
	return nil
}