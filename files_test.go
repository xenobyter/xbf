package main

import (
	"reflect"
	"testing"
)

func TestSortFiles(t *testing.T) {
	input := []FileInfo{
		{Name: "zebra.txt", IsDir: false},
		{Name: "Beta_dir", IsDir: true},
		{Name: "apple.txt", IsDir: false},
		{Name: "alpha_dir", IsDir: true},
		{Name: "B.txt", IsDir: false},
		{Name: "a.txt", IsDir: false},
	}

	expected := []FileInfo{
		{Name: "alpha_dir", IsDir: true},
		{Name: "Beta_dir", IsDir: true},
		{Name: "a.txt", IsDir: false},
		{Name: "apple.txt", IsDir: false},
		{Name: "B.txt", IsDir: false},
		{Name: "zebra.txt", IsDir: false},
	}

	sortFiles(input)

	for i := range expected {
		if !reflect.DeepEqual(input[i], expected[i]) {
			t.Errorf("at index %d: expected %+v, got %+v", i, expected[i], input[i])
		}
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{size: 0, want: "0 B"},
		{size: 1023, want: "1023 B"},
		{size: 1024, want: "1.0 KB"},
		{size: 1024 * 1024, want: "1.0 MB"},
	}

	for _, test := range tests {
		if got := formatFileSize(test.size); got != test.want {
			t.Errorf("formatFileSize(%d) = %q, want %q", test.size, got, test.want)
		}
	}
}
