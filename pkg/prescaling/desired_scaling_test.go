package prescaling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDesiredScaling(t *testing.T) {

	testCases := []struct {
		b              bool
		multiplier     int
		desiredReplica int
		expected       int
		currentReplica int32
		name           string
	}{
		{
			name:           "Test - In Time Desired > Current",
			b:              true,
			desiredReplica: 20,
			currentReplica: 10,
			multiplier:     1,
			expected:       11,
		},
		{
			name:           "Test - In Time and Multplier 0",
			b:              true,
			desiredReplica: 20,
			currentReplica: 10,
			multiplier:     0,
			expected:       11,
		},
		{
			name:           "Test - Out Time and Multplier 0",
			b:              false,
			desiredReplica: 20,
			currentReplica: 10,
			multiplier:     0,
			expected:       0,
		},
		{
			name:           "Test - Desired equal Current but multiplier 2",
			b:              true,
			desiredReplica: 10,
			currentReplica: 10,
			multiplier:     2,
			expected:       11,
		},
		{
			name:           "Test - Desired equal Current",
			b:              true,
			desiredReplica: 10,
			currentReplica: 10,
			multiplier:     0,
			expected:       10,
		},
		{
			name:           "Test - Current > Desired",
			b:              true,
			desiredReplica: 10,
			currentReplica: 20,
			multiplier:     0,
			expected:       9,
		},
		{
			name:           "Test - Current > Desired with multiplier",
			b:              true,
			desiredReplica: 10,
			currentReplica: 30,
			multiplier:     2,
			expected:       9,
		},
		{
			name:           "Test - Current  > Desired but less 10% ",
			b:              true,
			desiredReplica: 50,
			currentReplica: 51,
			multiplier:     0,
			expected:       10,
		},
	}

	for _, testCase := range testCases {
		e := DesiredScaling(testCase.b, testCase.multiplier, testCase.desiredReplica, testCase.currentReplica)
		assert.Equal(t, e, testCase.expected, testCase.name)
	}
}
