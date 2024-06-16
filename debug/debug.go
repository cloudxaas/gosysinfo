package cxsysinfodebug

import (
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
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
	buf = append(buf, "CPU: "...) //cpu id
	buf = strconv.AppendInt(buf, int64(cxcputhread.CPUThread), 10)
	buf = append(buf, " \tAl: "...) //allocation
	buf = strconv.AppendInt(buf, int64(m.Alloc), 10)
	buf = append(buf, " B\tTA: "...) //total alloc
	buf = strconv.AppendInt(buf, int64(m.TotalAlloc), 10)
	buf = append(buf, " B\tSys: "...) //sys memory
	buf = strconv.AppendInt(buf, int64(m.Sys), 10)
	buf = append(buf, " B\tGCNo: "...) //number of gc
	buf = strconv.AppendInt(buf, int64(m.NumGC), 10)
	buf = append(buf, "\tHpSys: "...) //heap of sys
	buf = strconv.AppendInt(buf, int64(m.HeapSys), 10)
	buf = append(buf, " B\tHpUse: "...) //heap in use
	buf = strconv.AppendInt(buf, int64(m.HeapInuse), 10)
	buf = append(buf, " B\tHpObjs: "...) //heap objs
	buf = strconv.AppendInt(buf, int64(m.HeapObjects), 10)
	buf = append(buf, "\tGoNo: "...)//num of goroutine
	buf = strconv.AppendInt(buf, int64(runtime.NumGoroutine()), 10)
	buf = append(buf, "\tFD: "...) //file descriptor opened
	buf = strconv.AppendInt(buf, int64(tracker.OpenDescriptors), 10)
	buf = append(buf, "\tGC: "...) //garbage collection time
	buf = cxfmtreadable.FormatDuration(buf, stwPause)
	buf = append(buf, '\n')
	os.Stdout.Write(buf)
}

// The methods to handle the open/close of file descriptors are not shown here.
// These methods should update tracker.OpenDescriptors appropriately.
