package message

import (
	"encoding/json"

	"github.com/lirm/aeron-go/aeron/atomic"
)

// Codec handles message serialization for Aeron
type Codec struct{}

// NewCodec creates a new Codec instance
func NewCodec() *Codec {
	return &Codec{}
}

// Encode serializes a Message to bytes for Aeron
func (c *Codec) Encode(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

// Decode deserializes bytes from Aeron to a Message
func (c *Codec) Decode(buffer *atomic.Buffer, offset, length int32) (*Message, error) {
	data := make([]byte, length)
	buffer.GetBytes(offset, data)

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ToBuffer converts encoded message to Aeron buffer
func (c *Codec) ToBuffer(msg *Message) (*atomic.Buffer, int32, error) {
	data, err := c.Encode(msg)
	if err != nil {
		return nil, 0, err
	}
	buffer := atomic.MakeBuffer(data)
	return buffer, int32(len(data)), nil
}
