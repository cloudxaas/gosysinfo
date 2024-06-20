package cxsysinfodebug

import (
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
	cxcputhread "github.com/cloudxaas/gocpu/thread"
)

// Static buffer to prevent allocations
var buf [1024]byte
var stwPause time.Duration

// FileDescriptorTracker is a struct to track the number of open file descriptors.
type FileDescriptorTracker struct {
	OpenDescriptors int32
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
		runtime.GC()
		end := time.Now()
		stwPause = end.Sub(start)
		time.Sleep(period)
	}
}

func logStats(m *runtime.MemStats, tracker *FileDescriptorTracker) {
	idx := 0
	idx += copy(buf[idx:], "CPU: ")
	idx += copyInt(buf[idx:], int(cxcputhread.CPUThread))
	idx += copy(buf[idx:], " GC: ")
	idx += copyDuration(buf[idx:], time.Duration(m.PauseTotalNs))
	idx += copy(buf[idx:], " Al: ")
	idx += copyBytes(buf[idx:], m.Alloc)
	idx += copy(buf[idx:], " TA: ")
	idx += copyBytes(buf[idx:], m.TotalAlloc)
	idx += copy(buf[idx:], " Sys: ")
	idx += copyBytes(buf[idx:], m.Sys)
	idx += copy(buf[idx:], " GCNo: ")
	idx += copyInt(buf[idx:], int(m.NumGC))
	idx += copy(buf[idx:], " HpSys: ")
	idx += copyBytes(buf[idx:], m.HeapSys)
	idx += copy(buf[idx:], " HpUse: ")
	idx += copyBytes(buf[idx:], m.HeapInuse)
	idx += copy(buf[idx:], " HpObj: ")
	idx += copyCompactNumber(buf[idx:], int64(m.HeapObjects))
	idx += copy(buf[idx:], " GoNo: ")
	idx += copyInt(buf[idx:], runtime.NumGoroutine())
	idx += copy(buf[idx:], " FD: ")
	idx += copyInt(buf[idx:], int64(tracker.OpenDescriptors))
	buf[idx] = '\n'
	idx++
	os.Stdout.Write(buf[:idx])
}

// Helper functions to append various types to the buffer without causing allocations
func copyInt(dst []byte, num int) int {
	return copy(dst, strconv.Itoa(num))
}

func copyBytes(dst []byte, num uint64) int {
	return copy(dst, strconv.FormatUint(num, 10))
}

func copyDuration(dst []byte, d time.Duration) int {
	return copy(dst, d.String())
}

func copyCompactNumber(dst []byte, num int64) int {
	return copy(dst, strconv.FormatInt(num, 10))
}
