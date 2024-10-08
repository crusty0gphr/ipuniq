package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/bits-and-blooms/bitset"

	"github.com/ipuniq"
)

const numWorkers = 20

func main() {
	filePath := flag.String("file", "", "Path to the file to be processed")
	flag.Parse()
	if *filePath == "" {
		log.Fatalf("Error: file path must be provided using the -file argument")
	}

	startProfiling()

	startTime := time.Now()
	ipuniq.LogMemoryUsage("before function")

	bitSet := bitset.New(1 << 32)
	file, fileSize, err := ipuniq.OpenFile(*filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	chunks := ipuniq.SplitFileIntoChunks(file, fileSize, numWorkers)

	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)
		go func(i int, chunk ipuniq.ChunkMeta) {
			errProc := ipuniq.ProcessChunk(i, *filePath, chunk.StartOffset, chunk.EndOffset, bitSet, &wg)
			if errProc != nil {
				log.Printf("Error processing chunk %d: %v", i, errProc)
			}
		}(i, chunk)
	}
	wg.Wait()

	ipuniq.LogMemoryUsage("after function")
	log.Printf("Finished execution: %v", time.Since(startTime))
	log.Printf("Distinct IPs count: %v", bitSet.Count())
}

func startProfiling() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
