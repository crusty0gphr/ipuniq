package ipuniq

import (
	"testing"
)

func TestAdjustToNextNewline(t *testing.T) {
	filePath := "bucket/test_chunk"
	file, _, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	initialOffset := int64(5)
	adjustment, err := AdjustToNextNewline(file, initialOffset)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if adjustment < 0 {
		t.Fatalf("Expected adjustment to be greater than or equal to 0, got %d", adjustment)
	}

	adjustedOffset := initialOffset + adjustment
	_, err = file.Seek(adjustedOffset, 0)
	if err != nil {
		t.Fatalf("Error seeking to adjusted offset: %v", err)
	}

	buf := make([]byte, 1)
	_, err = file.Read(buf)
	if err != nil {
		t.Fatalf("Error reading from file: %v", err)
	}

	if buf[0] == '\n' {
		t.Fatalf("Expected non-newline character at adjusted offset %d, got newline", adjustedOffset)
	}
}
