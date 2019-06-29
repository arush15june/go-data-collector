package main

//               Client
// Client for transmitting data to gRPC server every 1s

import (
	"log"
	"os"
	"time"

	Collect "../lib/collect"
	Comm "../lib/communication"
	Logger "../lib/logger"
)

const (
	DefaultConnectionURL = "13.234.126.213:1337"
)

var (
	logger     *log.Logger
	connection *Comm.DataAPIConnection
)

func main() {
	// Initialize collector logger and server connections.
	logger = Logger.NewLogger("client")

	connection := Comm.NewConnection(DefaultConnectionURL, logger)
	connection.Dial()

	for {
		var (
			systemInfo = Collect.GetSystemInfo()
		)
		resp := connection.SendSystemInfo(systemInfo)
		logger.Println("Response: ", resp)
		time.Sleep(1000 * time.Millisecond)
	}

	os.Exit(0)
}
