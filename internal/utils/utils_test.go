package utils

import (
	"testing"
)

func TestShouldIgnoreFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"image.jpg", true},
		{"video.mp4", true},
		{"document.pdf", false},
		{"archive.zip", false},
		{"IMAGE.JPG", true}, // Case insensitive check
		{"script.sh", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if got := ShouldIgnoreFile(tt.filename); got != tt.want {
				t.Errorf("ShouldIgnoreFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestCalculateSHA256(t *testing.T) {
	data := []byte("hello world")
	// SHA256 of "hello world"
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	if got := CalculateSHA256(data); got != expected {
		t.Errorf("CalculateSHA256() = %v, want %v", got, expected)
	}
}
