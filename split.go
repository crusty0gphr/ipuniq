package ipuniq

import (
	"bufio"
	"log"
	"os"
)

// ChunkMeta defines the start and end offsets for a file chunk
type ChunkMeta struct {
	StartOffset int64
	EndOffset   int64
}

func OpenFile(path string) (*os.File, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		file.Close()
		return nil, 0, err
	}

	return file, fileInfo.Size(), nil
}

// SplitFileIntoChunks divides the file into equal-sized chunks for concurrent processing
// Each chunk ends at a newline boundary to ensure no lines are split between chunks
func SplitFileIntoChunks(file *os.File, fileSize int64, numWorkers int) []ChunkMeta {
	chunkSize := fileSize / int64(numWorkers)
	var chunks []ChunkMeta

	startOffset := int64(0)
	for i := 0; i < numWorkers; i++ {
		endOffset := startOffset + chunkSize

		// Adjust to the next newline, exclude the last line
		if i != numWorkers-1 {
			adjustment, err := AdjustToNextNewline(file, endOffset)
			if err != nil {
				log.Fatalf("Error adjusting to newline at chunk boundary: %v", err)
			}
			endOffset += adjustment
		} else {
			endOffset = fileSize
		}

		chunks = append(chunks, ChunkMeta{startOffset, endOffset})

		// Move startOffset
		startOffset = endOffset
	}

	return chunks
}

// AdjustToNextNewline moves the offset to the start of the next line to ensure a chunk ends at a newline
// This avoids splitting lines between chunks
func AdjustToNextNewline(file *os.File, offset int64) (int64, error) {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadBytes('\n')
	if err != nil && err.Error() != "EOF" {
		return 0, err
	}

	return int64(len(line)), nil
}
