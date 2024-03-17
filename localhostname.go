package cxsysinfo

import (
	"os"
)

var (
	LocalHostname string
)
	
func init() {
	LocalHostname, _ = os.Hostname()	
}
