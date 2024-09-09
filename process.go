package ipuniq

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

const bufferSize = 2 * 1024 * 1024 // 2 MB buffer

// ProcessChunk handles processing of a file chunk. Reads lines of the chunk withing the offsets
func ProcessChunk(id int, path string, startOffset, endOffset int64, set *Set, wg *sync.WaitGroup) error {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("worker %d: error opening file: %v", id, err)
	}
	defer file.Close()

	// Move file pointer to the specified start offset
	if _, err = file.Seek(startOffset, 0); err != nil {
		return fmt.Errorf("worker %d: unable to seek to startOffset: %v", id, err)
	}

	if err = processLinesInChunk(file, startOffset, endOffset, set, id); err != nil {
		return err
	}

	return nil
}

// processLinesInChunk reads lines from the file (withing the offsets) and processes them
// Checks for a valid IPv4 address and adds valid lines to the Set
func processLinesInChunk(file *os.File, startOffset, endOffset int64, set *Set, id int) error {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, bufferSize), bufferSize)

	curOffset := startOffset
	for scanner.Scan() {
		line := scanner.Bytes()
		curOffset += int64(len(line)) + 1 // +1 for the newline

		// Validate and store only valid IPv4 addresses
		ip := net.ParseIP(string(line))
		isValidIP := ip != nil && ip.To4() != nil
		if isValidIP {
			// Convert 16-bytes IP address to 4-bytes uint32
			set.Add(ipv4ToUint(ip))
		} else {
			// log.Printf("Worker %d: Invalid IPv4 address: %s\n", id, string(line))
			// TODO: add an invalid IP collector
		}

		// Stop if reached or exceeded the endOffset
		if curOffset >= endOffset {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("worker %d: error while reading file: %v", id, err)
	}

	return nil
}

func ipv4ToUint(ipv4 net.IP) uint32 {
	return binary.BigEndian.Uint32(ipv4.To4())
}
