package aeron

import (
	"context"
	"log/slog"
	"time"

	aeronlib "github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"

	"github.com/k-omotani/aeron-sample/internal/message"
)

// MessageHandler processes received messages
type MessageHandler func(msg *message.Message) error

// Subscriber wraps Aeron subscription for receiving messages
type Subscriber struct {
	subscription *aeronlib.Subscription
	codec        *message.Codec
	handler      MessageHandler
	logger       *slog.Logger
	idleStrategy idlestrategy.Idler
}

// NewSubscriber creates a subscriber on the given channel/stream
func NewSubscriber(
	aeron *aeronlib.Aeron,
	channel string,
	streamID int32,
	handler MessageHandler,
	logger *slog.Logger,
) (*Subscriber, error) {
	subscription, err := aeron.AddSubscription(channel, streamID)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		subscription: subscription,
		codec:        message.NewCodec(),
		handler:      handler,
		logger:       logger.With("component", "subscriber"),
		idleStrategy: idlestrategy.Sleeping{SleepFor: time.Millisecond},
	}, nil
}

// Start begins the polling loop in a goroutine
func (s *Subscriber) Start(ctx context.Context) {
	go s.pollLoop(ctx)
}

func (s *Subscriber) pollLoop(ctx context.Context) {
	s.logger.Info("subscriber poll loop started")

	fragmentHandler := func(buffer *atomic.Buffer, offset, length int32, header *logbuffer.Header) {
		msg, err := s.codec.Decode(buffer, offset, length)
		if err != nil {
			s.logger.Error("failed to decode message", "error", err)
			return
		}

		s.logger.Debug("received message",
			"type", msg.Type,
			"requestID", msg.RequestID,
			"timestamp", msg.Timestamp,
		)

		if err := s.handler(msg); err != nil {
			s.logger.Error("handler failed", "error", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("subscriber stopping")
			return
		default:
		}

		fragmentsRead := s.subscription.Poll(fragmentHandler, 10)
		if fragmentsRead == 0 {
			s.idleStrategy.Idle(0)
		}
	}
}

// Close releases the subscription resources
func (s *Subscriber) Close() error {
	return s.subscription.Close()
}
