package cxsysinfomem

import (
	"syscall"
	"os"
	"time"
	"strconv"
	"runtime"
)

func TotalPhysicalMemory() int {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	// If this is a 32-bit system, then these fields are
	// uint32 instead of uint64.
	// So we always convert to uint64 to match signature.
	return int(in.Totalram) * int(in.Unit)
}

func LogMemStatsPeriodically(period time.Duration) {
    var m runtime.MemStats
    var buf []byte
    for {
        runtime.ReadMemStats(&m)
        buf = append(buf[:0], "Alloc = "...)
        buf = strconv.AppendInt(buf, int64(m.Alloc), 10)
        buf = append(buf, " B\tTotalAlloc = "...)
        buf = strconv.AppendInt(buf, int64(m.TotalAlloc), 10)
        buf = append(buf, " B\tSys = "...)
        buf = strconv.AppendInt(buf, int64(m.Sys), 10)
        buf = append(buf, " B\tNumGC = "...)
        buf = strconv.AppendInt(buf, int64(m.NumGC), 10)
	buf = append(buf, " \tHeapSys = "...)
        buf = strconv.AppendInt(buf, int64(m.HeapSys), 10)
	buf = append(buf, " B\tHeapInuse = "...)
        buf = strconv.AppendInt(buf, int64(m.HeapInuse), 10)
	buf = append(buf, " B\tHeapObjects = "...)
        buf = strconv.AppendInt(buf, int64(m.HeapObjects), 10)
        buf = append(buf, '\n')
        os.Stdout.Write(buf)
        time.Sleep(period)
    }
}
