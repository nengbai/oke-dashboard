package main

import (
	"oke-creating/helper"
)

// https://github.com/oracle/oci-go-sdk/blob/master/example/example_containerengine_test.go
func main() {
	// var dlog DefaultSDKLogger
	// dlog.currentLoggingLevel = 2
	// dlog.debugLogger = log.New(os.Stderr, "DEBUG ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	// SetSDKLogger(dlog)
	helper.CreateOKE()
}
