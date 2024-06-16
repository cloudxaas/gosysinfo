
package cxsysinfodebug

import (
        "os"
        "runtime"
        "testing"
)

// Mock implementation for FileDescriptorTracker, adjust as necessary
type mockTracker struct {
        OpenDescriptors int32
}

func (mt *mockTracker) Increment() {
        mt.OpenDescriptors++
}

func (mt *mockTracker) Decrement() {
        if mt.OpenDescriptors > 0 {
                mt.OpenDescriptors--
        }
}

// BenchmarkLogStats tests the logStats function for allocations and time per operation.
func BenchmarkLogStats(b *testing.B) {
        // Mock setup
        tracker := &FileDescriptorTracker{}
        var m runtime.MemStats
        runtime.ReadMemStats(&m) // Get some initial real data to work with

        // Disable output during benchmarks to avoid measuring I/O time
        old := os.Stdout
        os.Stdout, _ = os.Open(os.DevNull)
        defer func() {
                os.Stdout = old
        }()

        b.ResetTimer() // Start timing after setup

        for i := 0; i < b.N; i++ {
                logStats(&m, tracker) // Call the function being benchmarked
        }

        b.StopTimer() // Stop timing before cleanup
}
