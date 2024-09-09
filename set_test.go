package ipuniq

import (
	"testing"
)

func TestBitMap_SetAndCount(t *testing.T) {
	set := NewSet()

	set.Add([]byte{127, 0, 0, 1})
	set.Add([]byte{192, 168, 0, 1})
	set.Add([]byte{127, 0, 0, 1})

	count := set.Count()
	if count != 2 {
		t.Errorf("Expected 2 distinct IPs, got %d", count)
	}

	set.Add([]byte{10, 0, 0, 1})
	set.Add([]byte{172, 16, 0, 1})

	count = set.Count()
	if count != 4 {
		t.Errorf("Expected 4 distinct IPs, got %d", count)
	}
}

func TestBytesToUint32(t *testing.T) {
	tests := []struct {
		input    []byte
		expected uint32
	}{
		{[]byte{127, 0, 0, 1}, 2130706433},   // 127.0.0.1 in uint32
		{[]byte{192, 168, 0, 1}, 3232235521}, // 192.168.0.1 in uint32
		{[]byte{10, 0, 0, 1}, 167772161},     // 10.0.0.1 in uint32
	}

	for _, tt := range tests {
		result := hashElement(tt.input)
		if result != tt.expected {
			t.Errorf("bytesToUint32(%v): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}
