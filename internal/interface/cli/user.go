package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/internal/usecase"
	useruc "github.com/Romasmi/s-shop-microservices/internal/usecase/user"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(resetPasswordCmd)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management",
}

var createUserCmd = &cobra.Command{
	Use:   "create [name] [email]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		email := args[1]

		handler := deps.GetHandler(usecase.UseCaseCreateUser)
		res, err := handler.Do(context.Background(), useruc.CreateUserInput{Name: name, Email: email})
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		u := res.(*user.User)
		fmt.Printf("User created: %v\n", u)
	},
}

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password [email]",
	Short: "Reset user password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		// Dummy implementation for example
		fmt.Printf("Password reset link sent to: %s\n", email)
	},
}
