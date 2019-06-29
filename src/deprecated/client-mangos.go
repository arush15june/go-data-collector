package main

//               Client
// Connect to the server at DefaultConnectionUrl,
// Collect system data using Collect.GetSystemInfo,
// Pack it as JSON,
// Transmit to the server every 1 second

import (
	"log"
	"os"
	"time"

	Collect "../lib/collect"
	Comm "../lib/communication"
	Logger "../lib/logger"
	Util "../lib/util"
)

const DefaultConnectionUrl = "tcp://127.0.0.0:1337"

var (
	logger     *log.Logger
	connection Comm.Connection
)

func main() {
	// Initialize collector logger and server connections.
	logger = Logger.NewLogger("client")

	connection := Comm.NewConnection(DefaultConnectionUrl, "PAIR", logger, false)
	connection.Run()

	for {
		var (
			systemInfo = Collect.GetSystemInfo()
			jsonString = Util.PackJSON(systemInfo)
		)

		connection.Send([]byte(jsonString))
		time.Sleep(1000 * time.Millisecond)
	}

	os.Exit(0)
}
