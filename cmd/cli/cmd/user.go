package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type User struct {
	Username string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
}

type UsersConfig struct {
	Users []User `yaml:"users"`
}

func NewUserAddCommand() *cobra.Command {
	var username, password, yamlFile string

	userAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new user",
		Long:  "Add a new user to the YAML configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			// Set default YAML file path if not provided
			if yamlFile == "" {
				yamlFile = "users.yaml"
			}

			// Check if YAML file exists, create if not
			if err := checkFileExists(yamlFile); err != nil {
				fmt.Printf("Error creating YAML file: %v\n", err)
				os.Exit(1)
			}

			// Read existing users
			var config UsersConfig
			data, err := os.ReadFile(yamlFile)
			if err != nil {
				fmt.Printf("Error reading YAML file: %v\n", err)
				os.Exit(1)
			}

			if err := yaml.Unmarshal(data, &config); err != nil {
				fmt.Printf("Error parsing YAML: %v\n", err)
				os.Exit(1)
			}

			// Check if user already exists
			for _, user := range config.Users {
				if user.Username == username {
					fmt.Printf("User '%s' already exists\n", username)
					os.Exit(1)
				}
			}

			// Hash the password (simple hash for now, in production use bcrypt)
			// For testing purposes, we'll use a simple hash
			passwordHash := hashPassword(password)

			// Add new user
			config.Users = append(config.Users, User{
				Username:     username,
				PasswordHash: passwordHash,
			})

			// Write back to YAML file
			updatedData, err := yaml.Marshal(&config)
			if err != nil {
				fmt.Printf("Error marshaling YAML: %v\n", err)
				os.Exit(1)
			}

			if err := os.WriteFile(yamlFile, updatedData, 0644); err != nil {
				fmt.Printf("Error writing YAML file: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("User '%s' created successfully\n", username)
		},
	}

	userAddCmd.Flags().StringVarP(&username, "username", "u", "", "Username for the new user (required)")
	userAddCmd.Flags().StringVarP(&password, "password", "p", "", "Password for the new user (required)")
	userAddCmd.Flags().StringVarP(&yamlFile, "file", "f", "", "YAML file path (default: users.yaml)")

	// Mark required flags
	_ = userAddCmd.MarkFlagRequired("username")
	_ = userAddCmd.MarkFlagRequired("password")

	return userAddCmd
}

func hashPassword(password string) string {
	// Simple hash function for testing
	// In production, use bcrypt or similar
	return fmt.Sprintf("hashed_%s", password)
}
