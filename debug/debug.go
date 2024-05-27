package cxsysinfodebug

import (
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
	cxcputhread "github.com/cloudxaas/gocpu/thread"
)

var stwPause time.Duration

// FileDescriptorTracker is a struct to track the number of open file descriptors.
type FileDescriptorTracker struct {
	OpenDescriptors int32
}


// FormatDuration formats a time.Duration into a human-readable string without heap allocations.
func FormatDuration(buf []byte, d time.Duration) []byte {
    switch {
    case d < time.Microsecond:
        buf = append(buf, strconv.Itoa(int(d.Nanoseconds()))...)
        buf = append(buf, "ns"...)
    case d < time.Millisecond:
        buf = append(buf, strconv.Itoa(int(d.Microseconds()))...)
        buf = append(buf, "µs"...)
    case d < time.Second:
        buf = append(buf, strconv.Itoa(int(d.Milliseconds()))...)
        buf = append(buf, "ms"...)
    case d < time.Minute:
        buf = append(buf, strconv.Itoa(int(d.Seconds()))...)
        buf = append(buf, "s"...)
    case d < time.Hour:
        buf = append(buf, strconv.Itoa(int(d.Minutes()))...)
        buf = append(buf, "m"...)
    default:
        buf = append(buf, strconv.Itoa(int(d.Hours()))...)
        buf = append(buf, "h"...)
    }
    return buf
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
	go recordPauseTime(period)
	for {
		runtime.ReadMemStats(&m)
		logStats(&m, tracker)
		time.Sleep(period)
	}
}

func recordPauseTime(period time.Duration) {
      debug.SetGCPercent(-1)

      for {
              start := time.Now()
              runtime.GC() // Trigger garbage collection
              end := time.Now()
              stwPause = end.Sub(start)
              time.Sleep(period) // Adjust the frequency of GC triggers
      }
}


func logStats(m *runtime.MemStats, tracker *FileDescriptorTracker) {
	buf := make([]byte, 0, 1024) // Preallocate buffer to avoid allocations
	buf = append(buf, "CPUThread = "...)
	buf = strconv.AppendInt(buf, int64(cxcputhread.CPUThread), 10)
	buf = append(buf, " \tAlloc = "...)
	buf = strconv.AppendInt(buf, int64(m.Alloc), 10)
	buf = append(buf, " B\tTotalAlloc = "...)
	buf = strconv.AppendInt(buf, int64(m.TotalAlloc), 10)
	buf = append(buf, " B\tSys = "...)
	buf = strconv.AppendInt(buf, int64(m.Sys), 10)
	buf = append(buf, " B\tNumGC = "...)
	buf = strconv.AppendInt(buf, int64(m.NumGC), 10)
	buf = append(buf, "\tHeapSys = "...)
	buf = strconv.AppendInt(buf, int64(m.HeapSys), 10)
	buf = append(buf, " B\tHeapUse = "...)
	buf = strconv.AppendInt(buf, int64(m.HeapInuse), 10)
	buf = append(buf, " B\tHeapObjs = "...)
	buf = strconv.AppendInt(buf, int64(m.HeapObjects), 10)
	buf = append(buf, "\tNumGo = "...)
	buf = strconv.AppendInt(buf, int64(runtime.NumGoroutine()), 10)
	buf = append(buf, "\tOpenFD = "...)
	buf = strconv.AppendInt(buf, int64(tracker.OpenDescriptors), 10)
  buf = append(buf, "\tGCNSTime = "...)
    buf = FormatDuration(buf, stwPause)
  buf = append(buf, '\n')

	os.Stdout.Write(buf)
}

// The methods to handle the open/close of file descriptors are not shown here.
// These methods should update tracker.OpenDescriptors appropriately.
