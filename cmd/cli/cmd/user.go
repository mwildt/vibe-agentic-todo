package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vibe-agentic/auth/user"
)

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

			// Create user repository
			repo := user.NewYAMLUserRepository(yamlFile)
			if err := repo.LoadUsers(); err != nil {
				fmt.Printf("Error loading users: %v\n", err)
				os.Exit(1)
			}

			// Create the user
			if err := repo.CreateUser(username, password); err != nil {
				fmt.Printf("Error creating user: %v\n", err)
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

func NewUserUpdateCommand() *cobra.Command {
	var username, password, yamlFile string

	userUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing user",
		Long:  "Update an existing user in the YAML configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			// Set default YAML file path if not provided
			if yamlFile == "" {
				yamlFile = "users.yaml"
			}

			// Create user repository
			repo := user.NewYAMLUserRepository(yamlFile)
			if err := repo.LoadUsers(); err != nil {
				fmt.Printf("Error loading users: %v\n", err)
				os.Exit(1)
			}

			// Update the user
			if err := repo.UpdateUser(username, password); err != nil {
				fmt.Printf("Error updating user: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("User '%s' updated successfully\n", username)
		},
	}

	userUpdateCmd.Flags().StringVarP(&username, "username", "u", "", "Username of the user to update (required)")
	userUpdateCmd.Flags().StringVarP(&password, "password", "p", "", "New password for the user (required)")
	userUpdateCmd.Flags().StringVarP(&yamlFile, "file", "f", "", "YAML file path (default: users.yaml)")

	// Mark required flags
	_ = userUpdateCmd.MarkFlagRequired("username")
	_ = userUpdateCmd.MarkFlagRequired("password")

	return userUpdateCmd
}

func NewUserDeleteCommand() *cobra.Command {
	var username, yamlFile string

	userDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Long:  "Delete a user from the YAML configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			// Set default YAML file path if not provided
			if yamlFile == "" {
				yamlFile = "users.yaml"
			}

			// Create user repository
			repo := user.NewYAMLUserRepository(yamlFile)
			if err := repo.LoadUsers(); err != nil {
				fmt.Printf("Error loading users: %v\n", err)
				os.Exit(1)
			}

			// Delete the user
			if err := repo.DeleteUser(username); err != nil {
				fmt.Printf("Error deleting user: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("User '%s' deleted successfully\n", username)
		},
	}

	userDeleteCmd.Flags().StringVarP(&username, "username", "u", "", "Username of the user to delete (required)")
	userDeleteCmd.Flags().StringVarP(&yamlFile, "file", "f", "", "YAML file path (default: users.yaml)")

	// Mark required flags
	_ = userDeleteCmd.MarkFlagRequired("username")

	return userDeleteCmd
}

func hashPassword(password string) string {
	// Simple hash function for testing
	// In production, use bcrypt or similar
	return fmt.Sprintf("hashed_%s", password)
}
