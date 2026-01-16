package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	aeronlib "github.com/lirm/aeron-go/aeron"

	"github.com/k-omotani/aeron-sample/internal/aeron"
	"github.com/k-omotani/aeron-sample/internal/handler"
	"github.com/k-omotani/aeron-sample/internal/logging"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse flags
	httpAddr := flag.String("addr", ":8080", "HTTP listen address")
	logLevel := flag.String("log-level", "debug", "Log level (debug, info, warn, error)")
	aeronDir := flag.String("aeron-dir", "/dev/shm/aeron", "Aeron media driver directory")
	channel := flag.String("channel", "", "Aeron channel (e.g., aeron:udp?endpoint=subscriber-driver:40123)")
	streamID := flag.Int("stream-id", 1001, "Aeron stream ID")
	flag.Parse()

	// Setup logging
	logCfg := logging.DefaultConfig()
	logCfg.Level = logging.ParseLevel(*logLevel)
	logger := logging.NewLogger(logCfg)

	// Use environment variable if flag not provided
	channelStr := *channel
	if channelStr == "" {
		channelStr = os.Getenv("CHANNEL")
	}
	if channelStr == "" {
		channelStr = aeron.DefaultPublisherConfig().Channel
	}

	logger.Info("starting publisher application",
		"addr", *httpAddr,
		"aeronDir", *aeronDir,
		"channel", channelStr,
		"streamID", *streamID,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration
	config := aeron.DefaultPublisherConfig()
	config.AeronDir = *aeronDir
	config.Channel = channelStr
	config.StreamID = int32(*streamID)

	// Initialize Aeron
	aeronCtx := aeronlib.NewContext()
	aeronCtx.AeronDir(config.AeronDir)
	aeronCtx.MediaDriverTimeout(config.MediaDriverTimeout)
	aeronCtx.ErrorHandler(func(err error) {
		logger.Error("aeron error", "error", err)
	})

	aeronClient, err := aeronlib.Connect(aeronCtx)
	if err != nil {
		return fmt.Errorf("failed to connect to Aeron: %w", err)
	}
	defer aeronClient.Close()

	logger.Info("connected to Aeron media driver")

	// Initialize publisher
	publisher, err := aeron.NewPublisher(
		aeronClient,
		config.Channel,
		config.StreamID,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create publisher: %w", err)
	}
	defer publisher.Close()

	// Setup HTTP handlers
	publishHandler := handler.NewPublishHandler(publisher, logger)
	healthHandler := handler.NewHealthHandler()

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/counter/increment", publishHandler.Increment)
	mux.HandleFunc("GET /health", healthHandler.Health)
	mux.HandleFunc("GET /ready", healthHandler.Ready)

	// Create HTTP server
	server := &http.Server{
		Addr:         *httpAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP server
	go func() {
		logger.Info("starting HTTP server", "addr", *httpAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	select {
	case <-sigChan:
		logger.Info("shutdown signal received")
	case <-ctx.Done():
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}

	logger.Info("publisher shutdown complete")
	return nil
}
