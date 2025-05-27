package slicecommon

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		wantNums  []string
		wantWords []string
	}{
		{
			name:      "EmptySlice",
			input:     []string{},
			wantNums:  []string{},
			wantWords: []string{},
		},
		{
			name:      "OnlyNumbers",
			input:     []string{"123", "456", "7890"},
			wantNums:  []string{"123", "456", "7890"},
			wantWords: []string{},
		},
		{
			name:      "OnlyWords",
			input:     []string{"hello", "world", "golang"},
			wantNums:  []string{},
			wantWords: []string{"hello", "world", "golang"},
		},
		{
			name:      "MixedInput",
			input:     []string{"123", "hello", "456", "world", "789"},
			wantNums:  []string{"123", "456", "789"},
			wantWords: []string{"hello", "world"},
		},
		{
			name:      "NumbersWithLetters",
			input:     []string{"123abc", "456", "abc123", "789"},
			wantNums:  []string{"456", "789"},
			wantWords: []string{"123abc", "abc123"},
		},
		{
			name:      "SpecialCharacters",
			input:     []string{"123", "@#$", "456", "word!", "789"},
			wantNums:  []string{"123", "456", "789"},
			wantWords: []string{"@#$", "word!"},
		},
		{
			name:      "EmptyStrings",
			input:     []string{"", "123", "", "hello"},
			wantNums:  []string{"123"},
			wantWords: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanInput := make([]string, len(tt.input))
			for i, s := range tt.input {
				cleanInput[i] = strings.TrimSpace(s)
			}

			gotNums, gotWords := Partition(cleanInput)

			cleanGotNums := make([]string, len(gotNums))
			for i, s := range gotNums {
				cleanGotNums[i] = strings.TrimSpace(s)
			}
			cleanGotWords := make([]string, len(gotWords))
			for i, s := range gotWords {
				cleanGotWords[i] = strings.TrimSpace(s)
			}

			if !reflect.DeepEqual(cleanGotNums, tt.wantNums) {
				t.Errorf("Partition() nums = %v, want %v (raw got: %v)", cleanGotNums, tt.wantNums, gotNums)
			}
			if !reflect.DeepEqual(cleanGotWords, tt.wantWords) {
				t.Errorf("Partition() words = %v, want %v (raw got: %v)", cleanGotWords, tt.wantWords, gotWords)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		list     []T
		target   T
		expected bool
	}

	stringTests := []testCase[string]{
		{name: "string present", list: []string{"apple", "banana", "orange"}, target: "banana", expected: true},
		{name: "string absent", list: []string{"apple", "banana", "orange"}, target: "grape", expected: false},
		{name: "empty string slice", list: []string{}, target: "apple", expected: false},
	}

	intTests := []testCase[int]{
		{name: "int present", list: []int{1, 2, 3, 4}, target: 3, expected: true},
		{name: "int absent", list: []int{1, 2, 3, 4}, target: 5, expected: false},
	}

	boolTests := []testCase[bool]{
		{name: "bool true present", list: []bool{true, false}, target: true, expected: true},
		{name: "bool false absent", list: []bool{true}, target: false, expected: false},
	}

	t.Run("strings", func(t *testing.T) {
		for _, tt := range stringTests {
			t.Run(tt.name, func(t *testing.T) {
				result := Contains(tt.list, tt.target)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("integers", func(t *testing.T) {
		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				result := Contains(tt.list, tt.target)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("booleans", func(t *testing.T) {
		for _, tt := range boolTests {
			t.Run(tt.name, func(t *testing.T) {
				result := Contains(tt.list, tt.target)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}
