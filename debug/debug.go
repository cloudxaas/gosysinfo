package cxsysinfodebug

import (
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"
	cxcputhread "github.com/cloudxaas/gocpu/thread"
	cxfmtreadable "github.com/cloudxaas/gofmt/readable"
)

var stwPause time.Duration

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

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 1024)
	},
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
		stwPause = time.Since(start)
		time.Sleep(period) // Adjust the frequency of GC triggers
	}
}

func logStats(m *runtime.MemStats, tracker *FileDescriptorTracker) {
	buf := bufferPool.Get().([]byte)
	buf = buf[:0] // Reset slice length

	buf = append(buf, "CPU: "...) // CPU id
	buf = strconv.AppendInt(buf, int64(cxcputhread.CPUThread), 10)
	buf = append(buf, "  GC: "...) // garbage collection time
	buf = cxfmtreadable.FormatDuration(buf, time.Duration(m.PauseTotalNs)) // Convert uint64 to time.Duration
	buf = append(buf, "  Al: "...) // allocation
	buf = cxfmtreadable.AppendBytes(buf, uint64(m.Alloc))
	buf = append(buf, "  TA: "...) // total alloc
	buf = cxfmtreadable.AppendBytes(buf, uint64(m.TotalAlloc))
	buf = append(buf, "  Sys: "...) // sys memory
	buf = cxfmtreadable.AppendBytes(buf, uint64(m.Sys))
	buf = append(buf, "  GCNo: "...) // number of GC
	buf = strconv.AppendInt(buf, int64(m.NumGC), 10)
	buf = append(buf, "  HpSys: "...) // heap sys
	buf = cxfmtreadable.AppendBytes(buf, uint64(m.HeapSys))
	buf = append(buf, "  HpUse: "...) // heap in use
	buf = cxfmtreadable.AppendBytes(buf, uint64(m.HeapInuse))
	buf = append(buf, "  HpObj: "...) // heap objects
	buf = cxfmtreadable.FormatNumberCompact(int64(m.HeapObjects), buf)
	buf = append(buf, "  GoNo: "...) // number of goroutines
	buf = strconv.AppendInt(buf, int64(runtime.NumGoroutine()), 10)
	buf = append(buf, "  FD: "...) // file descriptors opened
	buf = strconv.AppendInt(buf, int64(tracker.OpenDescriptors), 10)
	buf = append(buf, '\n')
	os.Stdout.Write(buf)
	bufferPool.Put(buf)
}

func appendWithLabel(buf []byte, label string, value int64) []byte {
	buf = append(buf, label...)
	return strconv.AppendInt(buf, value, 10)
}
