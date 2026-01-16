package aeron

import (
	"context"
	"errors"
	"log/slog"
	"time"

	aeronlib "github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"

	"github.com/k-omotani/aeron-sample/internal/message"
)

var (
	ErrNotConnected  = errors.New("publication not connected")
	ErrBackPressured = errors.New("publication back pressured")
	ErrOfferFailed   = errors.New("offer failed")
)

// Publisher wraps Aeron publication for sending messages
type Publisher struct {
	publication *aeronlib.Publication
	codec       *message.Codec
	logger      *slog.Logger
}

// NewPublisher creates a publisher on the given channel/stream
func NewPublisher(aeron *aeronlib.Aeron, channel string, streamID int32, logger *slog.Logger) (*Publisher, error) {
	publication, err := aeron.AddPublication(channel, streamID)
	if err != nil {
		return nil, err
	}

	// Wait for publication to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for !publication.IsConnected() {
		select {
		case <-ctx.Done():
			return nil, ErrNotConnected
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	return &Publisher{
		publication: publication,
		codec:       message.NewCodec(),
		logger:      logger.With("component", "publisher"),
	}, nil
}

// Publish sends a message through Aeron
func (p *Publisher) Publish(ctx context.Context, msg *message.Message) error {
	buffer, length, err := p.codec.ToBuffer(msg)
	if err != nil {
		return err
	}

	return p.offer(ctx, buffer, length)
}

func (p *Publisher) offer(ctx context.Context, buffer *atomic.Buffer, length int32) error {
	maxRetries := 100
	retries := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		result := p.publication.Offer(buffer, 0, length, nil)

		switch {
		case result == aeronlib.NotConnected:
			p.logger.Warn("publication not connected, retrying")
			time.Sleep(100 * time.Millisecond)
		case result == aeronlib.BackPressured:
			p.logger.Debug("back pressured, retrying")
			time.Sleep(10 * time.Millisecond)
		case result < 0:
			retries++
			if retries > maxRetries {
				return ErrOfferFailed
			}
			time.Sleep(10 * time.Millisecond)
		default:
			p.logger.Debug("message published", "position", result)
			return nil
		}
	}
}

// Close releases the publication resources
func (p *Publisher) Close() error {
	return p.publication.Close()
}
