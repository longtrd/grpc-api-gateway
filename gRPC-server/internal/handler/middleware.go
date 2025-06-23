package handler

import (
	"context"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/longtrd/grpc-api-gateway/gRPC-server/internal/domain"
)

// LoggingInterceptor logs gRPC calls with duration and status
func LoggingInterceptor(logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Log the incoming request
		logger.Info("gRPC request started",
			"method", info.FullMethod,
			"start_time", start.Format(time.RFC3339),
		)

		// Call the handler
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		// Log the response
		logger.Info("gRPC request completed",
			"method", info.FullMethod,
			"duration_ms", duration.Milliseconds(),
			"status", statusCode.String(),
			"error", err,
		)

		return resp, err
	}
}

// RecoveryInterceptor recovers from panics and returns proper gRPC errors
func RecoveryInterceptor(logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				logger.Error("panic recovered in gRPC handler",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)

				// Return internal server error
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// ValidationInterceptor validates request context and basic requirements
func ValidationInterceptor(logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check if context is already cancelled
		if ctx.Err() != nil {
			logger.Warn("request context cancelled before processing",
				"method", info.FullMethod,
				"error", ctx.Err(),
			)
			return nil, status.Errorf(codes.Canceled, "request cancelled")
		}

		// Validate request is not nil
		if req == nil {
			logger.Warn("received nil request",
				"method", info.FullMethod,
			)
			return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
		}

		return handler(ctx, req)
	}
}

// MetricsInterceptor collects metrics for gRPC calls (placeholder)
func MetricsInterceptor(logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		// TODO: Implement actual metrics collection (Prometheus, etc.)
		logger.Debug("metrics collected",
			"method", info.FullMethod,
			"duration_ms", duration.Milliseconds(),
			"status", statusCode.String(),
		)

		return resp, err
	}
}

// AuthInterceptor handles authentication (placeholder for when auth is implemented)
func AuthInterceptor(logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for health checks and reflection
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// TODO: Implement actual authentication logic
		// For now, just log that auth would happen here
		logger.Debug("authentication check",
			"method", info.FullMethod,
			"note", "authentication not implemented yet",
		)

		return handler(ctx, req)
	}
}

// RateLimitInterceptor implements rate limiting (placeholder)
func RateLimitInterceptor(logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Implement actual rate limiting
		// For now, just log that rate limiting would happen here
		logger.Debug("rate limit check",
			"method", info.FullMethod,
			"note", "rate limiting not implemented yet",
		)

		return handler(ctx, req)
	}
}

// TimeoutInterceptor enforces timeout on gRPC calls
func TimeoutInterceptor(timeout time.Duration, logger domain.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Channel to receive result
		resultChan := make(chan result, 1)

		// Run handler in goroutine
		go func() {
			resp, err := handler(timeoutCtx, req)
			resultChan <- result{resp: resp, err: err}
		}()

		// Wait for either completion or timeout
		select {
		case res := <-resultChan:
			return res.resp, res.err
		case <-timeoutCtx.Done():
			logger.Warn("request timeout",
				"method", info.FullMethod,
				"timeout", timeout,
			)
			return nil, status.Errorf(codes.DeadlineExceeded, "request timeout")
		}
	}
}

// result is a helper struct for timeout interceptor
type result struct {
	resp interface{}
	err  error
}

// isPublicMethod checks if a method should skip authentication
func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}

	return false
}

// ChainInterceptors chains multiple interceptors together
func ChainInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Build the chain of interceptors
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := chain
			chain = func(ctx context.Context, req interface{}) (interface{}, error) {
				return interceptor(ctx, req, info, next)
			}
		}

		return chain(ctx, req)
	}
}
