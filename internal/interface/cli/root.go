package cli

import (
	"github.com/Romasmi/s-shop-microservices/internal/config"
	"github.com/Romasmi/s-shop-microservices/internal/usecase"
	"github.com/spf13/cobra"
)

// AppDependencies is an interface to avoid circular dependency if we use app.App
type AppDependencies interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
	GetConfig() *config.Config
}

var (
	deps AppDependencies
)

func SetApp(d AppDependencies) {
	deps = d
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "App CLI",
	Long:  `A command line interface for the microservice template.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands here
}
