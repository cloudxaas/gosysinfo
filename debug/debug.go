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
	idx += copyUint16(buf[idx:], cxcputhread.CPUThread)
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
	idx += copyInt(buf[idx:], int(m.HeapObjects))
	idx += copy(buf[idx:], " GoNo: ")
	idx += copyInt(buf[idx:], runtime.NumGoroutine())
	idx += copy(buf[idx:], " FD: ")
	idx += copyInt(buf[idx:], int(tracker.OpenDescriptors))
	buf[idx] = '\n'
	idx++
	os.Stdout.Write(buf[:idx])
}

// Helper functions to append various types to the buffer without causing allocations
func copyInt(dst []byte, num int) int {
	n := strconv.Itoa(num)
	return copy(dst, n)
}

func copyUint16(dst []byte, num uint16) int {
	n := strconv.Itoa(int(num))
	return copy(dst, n)
}

func copyBytes(dst []byte, num uint64) int {
	n := strconv.FormatUint(num, 10)
	return copy(dst, n)
}

func copyDuration(dst []byte, d time.Duration) int {
	n := d.String()
	return copy(dst, n)
}
