package benchmark

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/tomotakashimizu/performance-testing-toolkit/benchmarks/serializer/internal/models"
	"github.com/tomotakashimizu/performance-testing-toolkit/benchmarks/serializer/internal/serializers"
)

// Runner handles the execution of serialization benchmarks
type Runner struct {
	users       []models.User
	serializers []serializers.Serializer
}

// NewRunner creates a new benchmark runner
func NewRunner() *Runner {
	return &Runner{}
}

// SetTestData sets the test data for benchmarking
func (r *Runner) SetTestData(users []models.User) {
	r.users = users
}

// AddSerializer adds a serializer to be benchmarked
func (r *Runner) AddSerializer(s serializers.Serializer) {
	r.serializers = append(r.serializers, s)
}

// RunBenchmarks executes benchmarks for all serializers
func (r *Runner) RunBenchmarks(iterations int) ([]serializers.SerializationResult, error) {
	if len(r.users) == 0 {
		return nil, fmt.Errorf("no test data provided")
	}

	results := make([]serializers.SerializationResult, 0, len(r.serializers))

	for _, ser := range r.serializers {
		fmt.Printf("Running benchmark for %s...\n", ser.Name())
		result, err := r.benchmarkSerializer(ser, iterations)
		if err != nil {
			return nil, fmt.Errorf("error benchmarking %s: %w", ser.Name(), err)
		}
		results = append(results, result)
	}

	return results, nil
}

// benchmarkSerializer runs benchmark for a single serializer
func (r *Runner) benchmarkSerializer(ser serializers.Serializer, iterations int) (serializers.SerializationResult, error) {
	result := serializers.SerializationResult{
		SerializerName: ser.Name(),
		MarshalTimes:   make([]int64, iterations),
		UnmarshalTimes: make([]int64, iterations),
	}

	// Use first user for initial serialization to get data size
	firstUser := r.users[0]
	data, err := ser.Marshal(firstUser)
	if err != nil {
		return result, fmt.Errorf("initial marshal failed: %w", err)
	}
	result.DataSize = len(data)

	// Run iterations
	for i := 0; i < iterations; i++ {
		marshalTime, unmarshalTime, err := r.measureSingleIteration(ser, firstUser)
		if err != nil {
			return result, fmt.Errorf("iteration %d failed: %w", i+1, err)
		}
		result.MarshalTimes[i] = marshalTime
		result.UnmarshalTimes[i] = unmarshalTime
	}

	// Calculate statistics
	result.MarshalAvgNs = calculateAverage(result.MarshalTimes)
	result.MarshalMedianNs = calculateMedian(result.MarshalTimes)
	result.UnmarshalAvgNs = calculateAverage(result.UnmarshalTimes)
	result.UnmarshalMedianNs = calculateMedian(result.UnmarshalTimes)

	return result, nil
}

// measureSingleIteration measures marshal and unmarshal time for a single iteration
func (r *Runner) measureSingleIteration(ser serializers.Serializer, user models.User) (marshalTime, unmarshalTime int64, err error) {
	// Measure marshal time
	start := time.Now()
	data, err := ser.Marshal(user)
	marshalTime = time.Since(start).Nanoseconds()
	if err != nil {
		return 0, 0, fmt.Errorf("marshal failed: %w", err)
	}

	// Measure unmarshal time
	start = time.Now()
	_, err = ser.Unmarshal(data)
	unmarshalTime = time.Since(start).Nanoseconds()
	if err != nil {
		return 0, 0, fmt.Errorf("unmarshal failed: %w", err)
	}

	return marshalTime, unmarshalTime, nil
}

// RunSymmetryTests checks how empty slices and maps are handled
func (r *Runner) RunSymmetryTests() ([]serializers.SymmetryResult, error) {
	results := make([]serializers.SymmetryResult, 0, len(r.serializers))

	for _, ser := range r.serializers {
		fmt.Printf("Running symmetry test for %s...\n", ser.Name())
		result := r.testSymmetry(ser)
		results = append(results, result)
	}

	return results, nil
}

// testSymmetry tests how empty/nil slices and maps are handled
func (r *Runner) testSymmetry(ser serializers.Serializer) serializers.SymmetryResult {
	result := serializers.SymmetryResult{
		SerializerName: ser.Name(),
	}

	// Test empty slices
	userWithEmptySlice := models.User{
		ID:       1,
		Name:     "Test",
		Email:    "test@example.com",
		Age:      25,
		IsActive: true,
		Tags:     []string{}, // Empty slice
		Profile: models.Profile{
			FirstName:   "Test",
			LastName:    "User",
			Bio:         "Test bio",
			Avatar:      "test.jpg",
			SocialLinks: []models.Link{}, // Empty slice
			Preferences: models.Preferences{
				Theme:         "light",
				Language:      "en",
				Notifications: map[string]bool{},
				Privacy: models.PrivacySettings{
					ProfilePublic: true,
					EmailVisible:  false,
					ShowActivity:  true,
				},
			},
		},
		Settings: models.Settings{
			Language: "en",
			TimeZone: "UTC",
			Features: []string{},
			Limits:   map[string]int{},
		},
		Metadata:  map[string]interface{}{},
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := ser.Marshal(userWithEmptySlice)
	if err != nil {
		result.Details += fmt.Sprintf("Empty slice marshal error: %v; ", err)
	} else {
		restored, err := ser.Unmarshal(data)
		if err != nil {
			result.Details += fmt.Sprintf("Empty slice unmarshal error: %v; ", err)
		} else {
			// Functional equivalence test (length = 0)
			tagsEqual := len(userWithEmptySlice.Tags) == 0 && len(restored.Tags) == 0
			socialLinksEqual := len(userWithEmptySlice.Profile.SocialLinks) == 0 && len(restored.Profile.SocialLinks) == 0
			result.EmptySlicesOK = tagsEqual && socialLinksEqual

			// Strict symmetry test (exact type preservation)
			tagsStrictEqual := reflect.DeepEqual(userWithEmptySlice.Tags, restored.Tags)
			socialLinksStrictEqual := reflect.DeepEqual(userWithEmptySlice.Profile.SocialLinks, restored.Profile.SocialLinks)
			result.StrictEmptySlicesOK = tagsStrictEqual && socialLinksStrictEqual

			if !result.EmptySlicesOK {
				if !tagsEqual {
					result.Details += fmt.Sprintf("Empty slice length mismatch Tags: original len=%d, restored len=%d; ",
						len(userWithEmptySlice.Tags), len(restored.Tags))
				}
				if !socialLinksEqual {
					result.Details += fmt.Sprintf("Empty slice length mismatch SocialLinks: original len=%d, restored len=%d; ",
						len(userWithEmptySlice.Profile.SocialLinks), len(restored.Profile.SocialLinks))
				}
			}

			if !result.StrictEmptySlicesOK {
				if !tagsStrictEqual {
					result.Details += fmt.Sprintf("Strict empty slice type mismatch Tags: original=%#v, restored=%#v; ",
						userWithEmptySlice.Tags, restored.Tags)
				}
				if !socialLinksStrictEqual {
					result.Details += fmt.Sprintf("Strict empty slice type mismatch SocialLinks: original=%#v, restored=%#v; ",
						userWithEmptySlice.Profile.SocialLinks, restored.Profile.SocialLinks)
				}
			}
		}
	}

	// Test nil slices
	userWithNilSlice := models.User{
		ID:       2,
		Name:     "Test2",
		Email:    "test2@example.com",
		Age:      30,
		IsActive: false,
		Tags:     nil, // Nil slice
		Profile: models.Profile{
			FirstName:   "Test2",
			LastName:    "User2",
			Bio:         "Test bio 2",
			Avatar:      "test2.jpg",
			SocialLinks: nil, // Nil slice
			Preferences: models.Preferences{
				Theme:         "dark",
				Language:      "ja",
				Notifications: nil,
				Privacy: models.PrivacySettings{
					ProfilePublic: false,
					EmailVisible:  true,
					ShowActivity:  false,
				},
			},
		},
		Settings: models.Settings{
			Language: "ja",
			TimeZone: "JST",
			Features: nil,
			Limits:   nil,
		},
		Metadata:  nil,
		CreatedAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err = ser.Marshal(userWithNilSlice)
	if err != nil {
		result.Details += fmt.Sprintf("Nil slice marshal error: %v; ", err)
	} else {
		restored, err := ser.Unmarshal(data)
		if err != nil {
			result.Details += fmt.Sprintf("Nil slice unmarshal error: %v; ", err)
		} else {
			// Functional nil handling test (nil or zero length)
			result.NilSlicesOK = (userWithNilSlice.Tags == nil && (restored.Tags == nil || len(restored.Tags) == 0)) &&
				(userWithNilSlice.Profile.SocialLinks == nil && (restored.Profile.SocialLinks == nil || len(restored.Profile.SocialLinks) == 0))

			// Strict nil preservation test (nil stays nil)
			tagsStrictNil := userWithNilSlice.Tags == nil && restored.Tags == nil
			socialLinksStrictNil := userWithNilSlice.Profile.SocialLinks == nil && restored.Profile.SocialLinks == nil
			result.StrictNilSlicesOK = tagsStrictNil && socialLinksStrictNil

			if !result.NilSlicesOK {
				result.Details += fmt.Sprintf("Nil slice handling: original=nil, restored=%v; ",
					restored.Tags)
			}

			if !result.StrictNilSlicesOK {
				if !tagsStrictNil {
					result.Details += fmt.Sprintf("Strict nil slice preservation Tags: original=nil, restored=%#v; ",
						restored.Tags)
				}
				if !socialLinksStrictNil {
					result.Details += fmt.Sprintf("Strict nil slice preservation SocialLinks: original=nil, restored=%#v; ",
						restored.Profile.SocialLinks)
				}
			}
		}
	}

	// Test empty maps
	userWithEmptyMap := models.User{
		ID:       3,
		Name:     "Test3",
		Email:    "test3@example.com",
		Metadata: map[string]interface{}{}, // Empty map
		Profile: models.Profile{
			Preferences: models.Preferences{
				Notifications: map[string]bool{}, // Empty map
			},
		},
	}

	data, err = ser.Marshal(userWithEmptyMap)
	if err != nil {
		result.Details += fmt.Sprintf("Empty map marshal error: %v; ", err)
	} else {
		restored, err := ser.Unmarshal(data)
		if err != nil {
			result.Details += fmt.Sprintf("Empty map unmarshal error: %v; ", err)
		} else {
			// Functional equivalence test (length = 0)
			metadataEqual := len(userWithEmptyMap.Metadata) == 0 && len(restored.Metadata) == 0
			notificationsEqual := len(userWithEmptyMap.Profile.Preferences.Notifications) == 0 && len(restored.Profile.Preferences.Notifications) == 0
			result.EmptyMapsOK = metadataEqual && notificationsEqual

			// Strict symmetry test (exact type preservation)
			metadataStrictEqual := reflect.DeepEqual(userWithEmptyMap.Metadata, restored.Metadata)
			notificationsStrictEqual := reflect.DeepEqual(userWithEmptyMap.Profile.Preferences.Notifications, restored.Profile.Preferences.Notifications)
			result.StrictEmptyMapsOK = metadataStrictEqual && notificationsStrictEqual

			if !result.EmptyMapsOK {
				if !metadataEqual {
					result.Details += fmt.Sprintf("Empty map length mismatch Metadata: original len=%d, restored len=%d; ",
						len(userWithEmptyMap.Metadata), len(restored.Metadata))
				}
				if !notificationsEqual {
					result.Details += fmt.Sprintf("Empty map length mismatch Notifications: original len=%d, restored len=%d; ",
						len(userWithEmptyMap.Profile.Preferences.Notifications), len(restored.Profile.Preferences.Notifications))
				}
			}

			if !result.StrictEmptyMapsOK {
				if !metadataStrictEqual {
					result.Details += fmt.Sprintf("Strict empty map type mismatch Metadata: original=%#v, restored=%#v; ",
						userWithEmptyMap.Metadata, restored.Metadata)
				}
				if !notificationsStrictEqual {
					result.Details += fmt.Sprintf("Strict empty map type mismatch Notifications: original=%#v, restored=%#v; ",
						userWithEmptyMap.Profile.Preferences.Notifications, restored.Profile.Preferences.Notifications)
				}
			}
		}
	}

	// Test nil maps
	userWithNilMap := models.User{
		ID:       4,
		Name:     "Test4",
		Email:    "test4@example.com",
		Metadata: nil, // Nil map
		Profile: models.Profile{
			Preferences: models.Preferences{
				Notifications: nil, // Nil map
			},
		},
	}

	data, err = ser.Marshal(userWithNilMap)
	if err != nil {
		result.Details += fmt.Sprintf("Nil map marshal error: %v; ", err)
	} else {
		restored, err := ser.Unmarshal(data)
		if err != nil {
			result.Details += fmt.Sprintf("Nil map unmarshal error: %v; ", err)
		} else {
			// Functional nil handling test (nil or zero length)
			result.NilMapsOK = (userWithNilMap.Metadata == nil && (restored.Metadata == nil || len(restored.Metadata) == 0)) &&
				(userWithNilMap.Profile.Preferences.Notifications == nil && (restored.Profile.Preferences.Notifications == nil || len(restored.Profile.Preferences.Notifications) == 0))

			// Strict nil preservation test (nil stays nil)
			metadataStrictNil := userWithNilMap.Metadata == nil && restored.Metadata == nil
			notificationsStrictNil := userWithNilMap.Profile.Preferences.Notifications == nil && restored.Profile.Preferences.Notifications == nil
			result.StrictNilMapsOK = metadataStrictNil && notificationsStrictNil

			if !result.NilMapsOK {
				result.Details += fmt.Sprintf("Nil map handling: original=nil, restored=%v; ",
					restored.Metadata)
			}

			if !result.StrictNilMapsOK {
				if !metadataStrictNil {
					result.Details += fmt.Sprintf("Strict nil map preservation Metadata: original=nil, restored=%#v; ",
						restored.Metadata)
				}
				if !notificationsStrictNil {
					result.Details += fmt.Sprintf("Strict nil map preservation Notifications: original=nil, restored=%#v; ",
						restored.Profile.Preferences.Notifications)
				}
			}
		}
	}

	if result.Details == "" {
		result.Details = "All tests passed"
	}

	return result
}

// calculateAverage calculates the average of a slice of int64 values
func calculateAverage(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	var sum int64
	for _, v := range values {
		sum += v
	}
	return sum / int64(len(values))
}

// calculateMedian calculates the median of a slice of int64 values
func calculateMedian(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int64, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}
