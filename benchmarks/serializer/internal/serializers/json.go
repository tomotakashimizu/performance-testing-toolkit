package serializers

import (
	"encoding/json"

	"github.com/tomotakashimizu/performance-testing-toolkit/benchmarks/serializer/internal/models"
)

// JSONSerializer implements Serializer interface for JSON
type JSONSerializer struct{}

// NewJSONSerializer creates a new JSONSerializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Name returns the name of the serializer
func (j *JSONSerializer) Name() string {
	return "JSON"
}

// Marshal serializes a User to JSON bytes
func (j *JSONSerializer) Marshal(user models.User) ([]byte, error) {
	return json.Marshal(user)
}

// Unmarshal deserializes JSON bytes to a User
func (j *JSONSerializer) Unmarshal(data []byte) (models.User, error) {
	var user models.User
	err := json.Unmarshal(data, &user)
	return user, err
}
