package commands

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const defaultPort = "50051"

var (
	port string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the gRPC server",
	Long:  "Start the gRPC User Service server on the specified port",
	Run:   runStart,
}

func init() {
	startCmd.Flags().StringVarP(&port, "port", "p", defaultPort, "Port to run the gRPC server on")
}

func runStart(cmd *cobra.Command, args []string) {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Starting gRPC server on port %s...", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
