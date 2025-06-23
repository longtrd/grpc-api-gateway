package commands

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/config"
	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/handler"
	pb "github.com/longtrd/grpc-api-gateway/gRPC-server/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
	// Load configuration
	cfg := config.Load()

	// Override port from command line if provided
	if port != defaultPort {
		cfg.Server.Port = port
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Build dependency container
	container, err := config.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to build container: %v", err)
	}
	defer func() {
		if err := container.Cleanup(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
	}()

	logger := container.GetLogger()
	logger.Info("Starting gRPC server setup")

	// Create listener
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to create listener", "address", address, "error", err)
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	// Create gRPC server with interceptors
	grpcServer := createGRPCServer(container)

	// Register services
	registerServices(grpcServer, container)

	// Setup graceful shutdown
	setupGracefulShutdown(grpcServer, container)

	logger.Info("Starting gRPC server", "address", address)
	log.Printf("Starting gRPC server on %s...", address)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve gRPC server", "error", err)
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

// createGRPCServer creates a gRPC server with all interceptors
func createGRPCServer(container *config.Container) *grpc.Server {
	cfg := container.GetConfig()
	logger := container.GetLogger()

	// Create interceptors
	interceptors := []grpc.UnaryServerInterceptor{
		handler.RecoveryInterceptor(logger),
		handler.LoggingInterceptor(logger),
		handler.ValidationInterceptor(logger),
	}

	// Add optional interceptors based on feature flags
	if cfg.Features.EnableMetrics {
		interceptors = append(interceptors, handler.MetricsInterceptor(logger))
	}

	if cfg.Features.EnableAuth {
		interceptors = append(interceptors, handler.AuthInterceptor(logger))
	}

	// Add timeout interceptor
	if cfg.Server.ReadTimeout > 0 {
		interceptors = append(interceptors, handler.TimeoutInterceptor(cfg.Server.ReadTimeout, logger))
	}

	// Server options
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.MaxRecvMsgSize(cfg.Server.MaxRecvSize),
		grpc.MaxSendMsgSize(cfg.Server.MaxSendSize),
	}

	return grpc.NewServer(opts...)
}

// registerServices registers all gRPC services
func registerServices(server *grpc.Server, container *config.Container) {
	logger := container.GetLogger()
	cfg := container.GetConfig()

	// Register health check service
	if cfg.Features.EnableHealthCheck {
		healthServer := health.NewServer()
		grpc_health_v1.RegisterHealthServer(server, healthServer)

		// Set serving status for services
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
		healthServer.SetServingStatus("user.UserService", grpc_health_v1.HealthCheckResponse_SERVING)

		logger.Info("Registered health check service")
	}

	// Enable reflection for development
	if cfg.Logger.Level == "debug" {
		reflection.Register(server)
		logger.Info("Registered reflection service")
	}

	// Register UserService
	userHandler := container.GetUserHandler()
	if userHandler != nil {
		pb.RegisterUserServiceServer(server, userHandler)
		logger.Info("Registered UserService")
	} else {
		logger.Error("UserHandler is nil, cannot register UserService")
	}

	logger.Info("All services registered successfully")
}

// setupGracefulShutdown handles graceful shutdown
func setupGracefulShutdown(server *grpc.Server, container *config.Container) {
	logger := container.GetLogger()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", "signal", sig.String())

		// Give requests time to finish
		done := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(done)
		}()

		// Wait for graceful shutdown or force stop after timeout
		select {
		case <-done:
			logger.Info("Server stopped gracefully")
		case <-time.After(30 * time.Second):
			logger.Warn("Force stopping server after timeout")
			server.Stop()
		}
	}()
}
