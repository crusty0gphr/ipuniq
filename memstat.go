package ipuniq

import (
	"log"
	"runtime"
	"time"
)

// MemoryStats holds different types of memory-related statistics.
type MemoryStats struct {
	HeapAlloc     float64 // Current heap memory allocation in MB
	Sys           float64 // Total system memory obtained from the OS in MB
	NumGC         uint32  // Number of completed garbage collection cycles
	NumGoroutines int     // Number of active goroutines
}

// GetMemoryStats returns detailed memory usage information.
func GetMemoryStats() MemoryStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return MemoryStats{
		HeapAlloc:     bytesToMB(memStats.HeapAlloc),
		Sys:           bytesToMB(memStats.Sys),
		NumGC:         memStats.NumGC,
		NumGoroutines: runtime.NumGoroutine(),
	}
}

func bytesToMB(b uint64) float64 {
	return float64(b) / (1024 * 1024)
}

func LogMemoryPeriodically(done chan bool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return // Stop the goroutine when done signal is received
		case <-ticker.C:
			LogMemoryUsage("periodic")
		}
	}
}

func LogMemoryUsage(phase string) {
	memoryStats := GetMemoryStats()
	log.Printf(
		"Memory usage %s: HeapAlloc: %.2f MB, Sys: %.2f MB, GC cycles: %d, Goroutines: %d",
		phase,
		memoryStats.HeapAlloc,
		memoryStats.Sys,
		memoryStats.NumGC,
		memoryStats.NumGoroutines,
	)
}
