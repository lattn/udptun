package udptun

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "udp => ", log.Ldate | log.Ltime)
