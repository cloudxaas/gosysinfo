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
	//cxfmtreadable "github.com/cloudxaas/gofmt/readable"
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

	buf = appendWithLabel(buf, "CPU: ", int64(cxcputhread.CPUThread))
	buf = appendWithLabel(buf, "  GC: ", int64(m.PauseTotalNs))
	buf = appendWithLabel(buf, "  Al: ", int64(m.Alloc))
	buf = appendWithLabel(buf, "  TA: ", int64(m.TotalAlloc))
	buf = appendWithLabel(buf, "  Sys: ", int64(m.Sys))
	buf = appendWithLabel(buf, "  GCNo: ", int64(m.NumGC))
	buf = appendWithLabel(buf, "  HpSys: ", int64(m.HeapSys))
	buf = appendWithLabel(buf, "  HpUse: ", int64(m.HeapInuse))
	buf = appendWithLabel(buf, "  HpObj: ", int64(m.HeapObjects))
	buf = appendWithLabel(buf, "  GoNo: ", int64(runtime.NumGoroutine()))
	buf = appendWithLabel(buf, "  FD: ", int64(tracker.OpenDescriptors))

	buf = append(buf, '\n')
	os.Stdout.Write(buf)

	bufferPool.Put(buf)
}

func appendWithLabel(buf []byte, label string, value int64) []byte {
	buf = append(buf, label...)
	return strconv.AppendInt(buf, value, 10)
}
