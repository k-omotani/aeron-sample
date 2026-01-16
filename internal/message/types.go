package message

import (
	"encoding/json"
	"time"
)

// MessageType identifies the type of message
type MessageType uint8

const (
	MessageTypeUnknown   MessageType = 0
	MessageTypeIncrement MessageType = 1
	MessageTypeReset     MessageType = 2
)

// Message represents the envelope for all Aeron messages
type Message struct {
	Type      MessageType `json:"type"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id"`
	Payload   []byte      `json:"payload,omitempty"`
}

// IncrementPayload contains increment-specific data
type IncrementPayload struct {
	Amount int64  `json:"amount"`
	Source string `json:"source"`
}

// NewIncrementMessage creates a new increment message
func NewIncrementMessage(requestID string, amount int64, source string) (*Message, error) {
	payload := IncrementPayload{
		Amount: amount,
		Source: source,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:      MessageTypeIncrement,
		Timestamp: time.Now().UnixNano(),
		RequestID: requestID,
		Payload:   payloadBytes,
	}, nil
}

// DecodeIncrementPayload extracts IncrementPayload from a Message
func (m *Message) DecodeIncrementPayload() (*IncrementPayload, error) {
	var payload IncrementPayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
