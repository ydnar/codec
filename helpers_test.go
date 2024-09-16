package codec

import (
	"reflect"
	"testing"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		size int
		want []int
	}{
		{"nil, 0", nil, 0, nil},
		{"nil, 1", nil, 1, make([]int, 1)},
		{"len=3, 0", make([]int, 3), 0, make([]int, 3)},
		{"len=3, 1", make([]int, 3), 1, make([]int, 3)},
		{"len=3, 2", make([]int, 3), 2, make([]int, 3)},
		{"len=3, 3", make([]int, 3), 3, make([]int, 3)},
		{"len=3, 10", make([]int, 3), 10, make([]int, 10)},
		{"preserves cap", make([]int, 3, 15), 10, make([]int, 10, 15)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Expand(&tt.s, tt.size)
			if len(tt.s) != len(tt.want) {
				t.Errorf("Expand: got %v, expected %v", tt.s, tt.want)
			}
			if cap(tt.s) != cap(tt.want) {
				t.Errorf("Expand: got cap %d, expected %v", cap(tt.s), cap(tt.want))
			}
			if !reflect.DeepEqual(tt.s, tt.want) {
				t.Errorf("Expand: got %v, expected %v", tt.s, tt.want)
			}
		})
	}
}
