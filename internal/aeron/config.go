package aeron

import (
	"time"
)

// Config holds Aeron-related configuration
type Config struct {
	// AeronDir is the media driver directory
	AeronDir string

	// Channel for pub/sub communication
	// Publisher: "aeron:udp?endpoint=subscriber-driver:40123"
	// Subscriber: "aeron:udp?endpoint=0.0.0.0:40123"
	Channel string

	// StreamID for the counter messages
	StreamID int32

	// Timeouts
	MediaDriverTimeout time.Duration
}

// DefaultPublisherConfig returns config for publisher (sends to subscriber)
func DefaultPublisherConfig() *Config {
	return &Config{
		AeronDir:           "/dev/shm/aeron",
		Channel:            "aeron:udp?endpoint=subscriber-driver:40123",
		StreamID:           1001,
		MediaDriverTimeout: 10 * time.Second,
	}
}

// DefaultSubscriberConfig returns config for subscriber (listens on UDP)
func DefaultSubscriberConfig() *Config {
	return &Config{
		AeronDir:           "/dev/shm/aeron",
		Channel:            "aeron:udp?endpoint=0.0.0.0:40123",
		StreamID:           1001,
		MediaDriverTimeout: 10 * time.Second,
	}
}
