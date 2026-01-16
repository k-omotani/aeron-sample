package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	aeronlib "github.com/lirm/aeron-go/aeron"

	"github.com/k-omotani/aeron-sample/internal/aeron"
	"github.com/k-omotani/aeron-sample/internal/counter"
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
	logLevel := flag.String("log-level", "debug", "Log level (debug, info, warn, error)")
	aeronDir := flag.String("aeron-dir", "/dev/shm/aeron", "Aeron media driver directory")
	channel := flag.String("channel", "", "Aeron channel (e.g., aeron:udp?endpoint=0.0.0.0:40123)")
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
		channelStr = aeron.DefaultSubscriberConfig().Channel
	}

	logger.Info("starting subscriber application",
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
	config := aeron.DefaultSubscriberConfig()
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

	// Initialize counter state
	counterState := counter.NewState()

	// Create message processor
	processor := counter.NewProcessor(counterState, logger)

	// Initialize subscriber
	subscriber, err := aeron.NewSubscriber(
		aeronClient,
		config.Channel,
		config.StreamID,
		processor.Handle,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create subscriber: %w", err)
	}
	defer subscriber.Close()

	// Start subscriber polling loop
	subscriber.Start(ctx)

	logger.Info("subscriber started, waiting for messages...")

	// Wait for shutdown signal
	select {
	case <-sigChan:
		logger.Info("shutdown signal received")
	case <-ctx.Done():
	}

	logger.Info("subscriber shutdown complete")
	return nil
}
