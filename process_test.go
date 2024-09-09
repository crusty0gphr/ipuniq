package ipuniq

import (
	"os"
	"sync"
	"testing"

	"github.com/bits-and-blooms/bitset"
)

func TestProcessChunk(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_chunk")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("127.0.0.1\n192.168.0.1\n10.0.0.1\n172.16.0.1\n127.0.0.1\n")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	bitmap := bitset.New(20)
	var wg sync.WaitGroup
	wg.Add(1)

	err = ProcessChunk(0, tmpFile.Name(), 0, int64(len(content)), bitmap, &wg)
	if err != nil {
		t.Fatalf("Error processing chunk: %v", err)
	}
	wg.Wait()

	count := bitmap.Count()
	if count != 4 {
		t.Errorf("Expected 4 distinct IPs, got %d", count)
	}
}
