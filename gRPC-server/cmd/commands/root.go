package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "gRPC User Service Server",
	Long:  "A gRPC server implementation for User Service with CRUD operations",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
}
