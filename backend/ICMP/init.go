package ICMP

import (
    "log"
)

var g_debug bool            // debug environment
var g_logger *log.Logger    // logger

func init() {
    g_debug = false
    g_logger = log.New(nil, "[ICMP] ", log.Lshortfile|log.Ldate)
}
