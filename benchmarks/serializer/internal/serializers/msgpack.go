package serializers

import (
	"github.com/tomotakashimizu/performance-testing-toolkit/benchmarks/serializer/internal/models"
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPackSerializer implements Serializer interface for MessagePack
type MsgPackSerializer struct{}

// NewMsgPackSerializer creates a new MsgPackSerializer
func NewMsgPackSerializer() *MsgPackSerializer {
	return &MsgPackSerializer{}
}

// Name returns the name of the serializer
func (m *MsgPackSerializer) Name() string {
	return "MessagePack"
}

// Marshal serializes a User to MessagePack bytes
func (m *MsgPackSerializer) Marshal(user models.User) ([]byte, error) {
	return msgpack.Marshal(user)
}

// Unmarshal deserializes MessagePack bytes to a User
func (m *MsgPackSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := msgpack.Unmarshal(data, &user)
	return user, err
}
