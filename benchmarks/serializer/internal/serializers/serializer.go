package serializers

import "github.com/tomotakashimizu/performance-testing-toolkit/benchmarks/serializer/internal/models"

// Serializer defines the interface for serialization operations
type Serializer interface {
	Name() string
	Marshal(user models.User) ([]byte, error)
	Unmarshal(data []byte) (models.User, error)
}

// SerializationResult contains the results of serialization benchmarks
type SerializationResult struct {
	SerializerName    string
	MarshalTimes      []int64 // nanoseconds
	UnmarshalTimes    []int64 // nanoseconds
	DataSize          int     // bytes
	MarshalAvgNs      int64
	MarshalMedianNs   int64
	UnmarshalAvgNs    int64
	UnmarshalMedianNs int64
}

// SymmetryResult contains the results of empty slice/map symmetry tests
type SymmetryResult struct {
	SerializerName string
	EmptySlicesOK  bool
	EmptyMapsOK    bool
	NilSlicesOK    bool
	NilMapsOK      bool
	Details        string
}
