package cxsysinfomem

import (
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// FileDescriptorTracker is a struct to track the number of open file descriptors.
type FileDescriptorTracker struct {
	OpenDescriptors int32
}

// TotalPhysicalMemory returns the total physical memory of the system.
func TotalPhysicalMemory() int {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	return int(in.Totalram) * int(in.Unit)
}

// LogMemStatsPeriodically logs memory and file descriptor stats periodically.
func LogMemStatsPeriodically(period time.Duration, tracker *FileDescriptorTracker) {
	var m runtime.MemStats
	for {
		runtime.ReadMemStats(&m)
		logStats(&m, tracker)
		time.Sleep(period)
	}
}

func logStats(m *runtime.MemStats, tracker *FileDescriptorTracker) {
	buf := make([]byte, 0, 1024) // Preallocate buffer to avoid allocations
	buf = strconv.AppendInt(buf, int64(m.Alloc), 10)
	buf = append(buf, " B\tTotalAlloc = "...)
	buf = strconv.AppendInt(buf, int64(m.TotalAlloc), 10)
	buf = append(buf, " B\tSys = "...)
	buf = strconv.AppendInt(buf, int64(m.Sys), 10)
	buf = append(buf, " B\tNumGC = "...)
	buf = strconv.AppendInt(buf, int64(m.NumGC), 10)
	buf = append(buf, "\tHeapSys = "...)
	buf = strconv.AppendInt(buf, int64(m.HeapSys), 10)
	buf = append(buf, " B\tHeapInuse = "...)
	buf = strconv.AppendInt(buf, int64(m.HeapInuse), 10)
	buf = append(buf, " B\tHeapObjects = "...)
	buf = strconv.AppendInt(buf, int64(m.HeapObjects), 10)
	buf = append(buf, "\tNumGoroutines = "...)
	buf = strconv.AppendInt(buf, int64(runtime.NumGoroutine()), 10)
	buf = append(buf, "\tOpenDescriptors = "...)
	buf = strconv.AppendInt(buf, int64(tracker.OpenDescriptors), 10)
	buf = append(buf, '\n')

	os.Stdout.Write(buf)
}

// The methods to handle the open/close of file descriptors are not shown here.
// These methods should update tracker.OpenDescriptors appropriately.
