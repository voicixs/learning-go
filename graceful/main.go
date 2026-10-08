// ### Why Graceful Shutdown Matters

// In a simple Go server using `http.ListenAndServe(":8080", mux)`, stopping the process (such as pressing `Ctrl+C` or when Docker / Kubernetes terminates a container) abruptly kills the program. 

// * **Without a graceful shutdown:** Any client in the middle of a request (such as a database write or file upload) gets disconnected immediately, leading to partial writes, data corruption, and dropped connections.
// * **With a graceful shutdown:** The server stops accepting **new** incoming connections, but keeps existing connections open until ongoing requests finish processing or until a configured timeout expires.

// ---

// ### Core Patterns Taught in the Tutorial

// 1. **Separating `createServer` and `runServer`**: Splitting server configuration from server lifecycle management makes the code modular and easily unit-testable.
// 2. **Non-blocking Server Startup**: Because `server.ListenAndServe()` is a blocking call, it is launched in a background goroutine so the main thread remains free to listen for termination signals.
// 3. **Buffered Error Channel (`chan error, 1`)**: Startup errors are passed back to the main thread through a buffered error channel. The expected `http.ErrServerClosed` error is filtered out so it isn't treated as a crash.
// 4. **Listening for OS Signals (`signal.Notify`)**: The app listens for OS signals like `os.Interrupt` (SIGINT / `Ctrl+C`) and `syscall.SIGTERM` (used by Kubernetes and Docker during pod termination) using a buffered `os.Signal` channel.
// 5. **Multiplexing with `select`**: A `select` block waits for either a startup error, an OS interrupt signal, or parent context cancellation (`ctx.Done()`).
// 6. **Timeout Context (`context.WithTimeout`)**: `server.Shutdown(shutdownCtx)` is passed a context with a deadline. If active requests do not finish within the allowed window, `server.Close()` forcefully tears down remaining connections to avoid hanging indefinitely.

// ---

// ### Complete Learning Code

// Here is the full, runnable implementation from the tutorial with comments explaining each step:


package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ============================================================================
// 1. SERVER CONFIGURATION & FACTORY
// ============================================================================
// createServer configures the router and HTTP server struct.
// Returning a pointer (*http.Server) allows callers to customize settings if needed.
func createServer() *http.Server {
	mux := http.NewServeMux()

	// Simulated slow endpoint: takes 8 seconds to complete
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Slow request started...")
		time.Sleep(8 * time.Second) // Simulates heavy processing or database queries
		fmt.Fprintf(w, "Slow request completed at %v\n", time.Now())
		log.Println("Slow request finished.")
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

// ============================================================================
// 2. SERVER LIFECYCLE & GRACEFUL SHUTDOWN
// ============================================================================
// runServer manages server startup, OS signal listening, and orderly shutdown.
func runServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	// Buffered error channel to catch startup errors without blocking the goroutine
	serverError := make(chan error, 1)

	// Run ListenAndServe in a goroutine because it is a blocking operation
	go func() {
		log.Println("Starting server on", server.Addr)
		err := server.ListenAndServe()

		// http.ErrServerClosed is the expected error when Shutdown is called
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
		close(serverError)
	}()

	// Buffered channel of size 1 to receive OS termination signals without blocking
	stop := make(chan os.Signal, 1)
	// Listen for Ctrl+C (os.Interrupt) and SIGTERM (Kubernetes/Docker shutdown)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Multiplex channel operations: wait for error, OS signal, or context cancellation
	select {
	case err := <-serverError:
		// Server failed immediately (e.g., port already in use)
		return fmt.Errorf("server startup failed: %w", err)

	case sig := <-stop:
		log.Printf("Shutdown signal received (%v). Initiating graceful shutdown...\n", sig)

	case <-ctx.Done():
		log.Println("Parent context cancelled. Initiating graceful shutdown...")
	}

	// Create a dedicated context with a deadline for the shutdown phase
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown stops accepting new connections and waits for active requests to complete
	if err := server.Shutdown(shutdownCtx); err != nil {
		// If graceful shutdown times out or fails, force-close remaining connections
		closeErr := server.Close()
		if closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return fmt.Errorf("graceful shutdown failed, forced close: %w", err)
	}

	log.Println("Server exited gracefully.")
	return nil
}

// ============================================================================
// MAIN EXECUTION
// ============================================================================

func main() {
	server := createServer()

	// Pass a root background context and a 10-second timeout for graceful shutdown
	timeout := 10 * time.Second
	if err := runServer(context.Background(), server, timeout); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
// ```

// ---

// ### How to Test Graceful Shutdown

// You can test this behavior in two terminals as demonstrated in the video:

// 1. **Terminal 1**: Start the server:
//    ```bash
//    go run main.go
//    ```
// 2. **Terminal 2**: Trigger the slow 8-second request:
//    ```bash
//    curl http://localhost:8080/slow
//    ```
// 3. **Terminal 1**: Immediately press `Ctrl+C` while the request is running.

// **Observation:**
// * The server prints `Shutdown signal received`.
// * Instead of exiting immediately and failing the curl command, the server waits until the 8-second request finishes, writes the response to Terminal 2, and then prints `Server exited gracefully` before exiting.

// *(If you set `shutdownTimeout` to `3 * time.Second` instead of 10 seconds, the context deadline will expire before 8 seconds, forcing the server to close immediately and returning `context deadline exceeded`.)*
