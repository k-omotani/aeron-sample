package counter

import (
	"log/slog"

	"github.com/k-omotani/aeron-sample/internal/message"
)

// Processor handles incoming messages and updates counter state
type Processor struct {
	state  *State
	logger *slog.Logger
}

// NewProcessor creates a new message processor
func NewProcessor(state *State, logger *slog.Logger) *Processor {
	return &Processor{
		state:  state,
		logger: logger.With("component", "processor"),
	}
}

// Handle processes a message and returns an error if processing fails
func (p *Processor) Handle(msg *message.Message) error {
	switch msg.Type {
	case message.MessageTypeIncrement:
		return p.handleIncrement(msg)
	case message.MessageTypeReset:
		return p.handleReset(msg)
	default:
		p.logger.Warn("unknown message type", "type", msg.Type, "requestID", msg.RequestID)
		return nil
	}
}

func (p *Processor) handleIncrement(msg *message.Message) error {
	payload, err := msg.DecodeIncrementPayload()
	if err != nil {
		p.logger.Error("failed to decode increment payload", "error", err, "requestID", msg.RequestID)
		return err
	}

	newValue := p.state.Increment(payload.Amount)

	p.logger.Info("counter incremented",
		"requestID", msg.RequestID,
		"amount", payload.Amount,
		"source", payload.Source,
		"newValue", newValue,
		"totalEvents", p.state.TotalEvents(),
	)

	return nil
}

func (p *Processor) handleReset(msg *message.Message) error {
	p.state.Reset()

	p.logger.Info("counter reset",
		"requestID", msg.RequestID,
	)

	return nil
}
