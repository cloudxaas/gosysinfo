package cxsysinfodebug

import (
    "os"
    "runtime"
    "strconv"
    "syscall"
    "time"
    "unsafe"
    cxcputhread "github.com/cloudxaas/gocpu/thread"
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

var buffer [1024]byte

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
    buf := buffer[:0]

    buf = appendLabel(buf, "CPU: ", int64(cxcputhread.CPUThread))
    buf = appendLabel(buf, "  GC: ", int64(m.PauseTotalNs))
    buf = appendLabel(buf, "  Al: ", int64(m.Alloc))
    buf = appendLabel(buf, "  TA: ", int64(m.TotalAlloc))
    buf = appendLabel(buf, "  Sys: ", int64(m.Sys))
    buf = appendLabel(buf, "  GCNo: ", int64(m.NumGC))
    buf = appendLabel(buf, "  HpSys: ", int64(m.HeapSys))
    buf = appendLabel(buf, "  HpUse: ", int64(m.HeapInuse))
    buf = appendLabel(buf, "  HpObj: ", int64(m.HeapObjects))
    buf = appendLabel(buf, "  GoNo: ", int64(runtime.NumGoroutine()))
    buf = appendLabel(buf, "  FD: ", int64(tracker.OpenDescriptors))
    buf = append(buf, '\n')

    os.Stdout.Write(buf)
}

func appendLabel(buf []byte, label string, value int64) []byte {
    buf = append(buf, label...)
    return strconv.AppendInt(buf, value, 10)
}

// FormatBytes formats a byte count as a human-readable string.
func FormatBytes(buf []byte, bytes uint64) []byte {
    const unit = 1024
    if bytes < unit {
        return strconv.AppendUint(buf, bytes, 10)
    }
    div, exp := uint64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    buf = strconv.AppendFloat(buf, float64(bytes)/float64(div), 'f', 1, 64)
    buf = append(buf, " "_BCDE[exp]...)
    return buf
}

// Using a string instead of an array to save on allocations
const _BCDE = " KMGTPE"

// FormatDuration formats a duration as a human-readable string.
func FormatDuration(buf []byte, d time.Duration) []byte {
    u := uint64(d)
    if u < uint64(time.Microsecond) {
        return strconv.AppendUint(buf, u, 10)
    }
    if u < uint64(time.Millisecond) {
        return append(strconv.AppendFloat(buf, float64(u)/float64(time.Microsecond), 'f', 3, 64), 'µ', 's')
    }
    if u < uint64(time.Second) {
        return append(strconv.AppendFloat(buf, float64(u)/float64(time.Millisecond), 'f', 3, 64), 'm', 's')
    }
    return append(strconv.AppendFloat(buf, float64(u)/float64(time.Second), 'f', 3, 64), 's')
}

// FormatNumberCompact formats a number in a compact form.
func FormatNumberCompact(buf []byte, n int64) []byte {
    abs := uint64(n)
    if n < 0 {
        abs = uint64(-n)
    }
    if abs < 1000 {
        if n < 0 {
            buf = append(buf, '-')
        }
        return strconv.AppendUint(buf, abs, 10)
    }
    if abs < 1000000 {
        return formatFloat(buf, float64(n)/1000, "K")
    }
    if abs < 1000000000 {
        return formatFloat(buf, float64(n)/1000000, "M")
    }
    return formatFloat(buf, float64(n)/1000000000, "B")
}

func formatFloat(buf []byte, f float64, suffix string) []byte {
    buf = strconv.AppendFloat(buf, f, 'f', 1, 64)
    l := len(buf)
    if buf[l-1] == '0' && buf[l-2] == '.' {
        buf = buf[:l-2]
    }
    return append(buf, suffix...)
}
