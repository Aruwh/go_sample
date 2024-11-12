package utils

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{
			name:     "Test case 1",
			slice1:   []int{1, 2, 3, 4, 5},
			slice2:   []int{3, 4, 5, 6, 7},
			expected: []int{1, 2},
		},
		{
			name:     "Test case 2",
			slice1:   []int{1, 2, 3, 4, 5},
			slice2:   []int{1, 2, 3, 4, 5},
			expected: []int{},
		},
		{
			name:     "Test case 3",
			slice1:   []int{1, 2, 3},
			slice2:   []int{4, 5, 6},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersection(tt.slice1, tt.slice2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Difference(%v, %v) = %v, expected %v", tt.slice1, tt.slice2, result, tt.expected)
			}
		})
	}
}
