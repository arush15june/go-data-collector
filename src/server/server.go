package main

//               Server
// Server for receiviing data from gRPC server every 1s

import (
	"log"
	"os"

	Comm "../lib/communication"
	Logger "../lib/logger"
)

const (
	DefaultConnectionURL = ":1337"
)

var (
	logger     *log.Logger
	connection *Comm.DataAPIConnection
)

// gRPC Server Handlers

func main() {
	// Initialize collector logger and server connections.
	logger = Logger.NewLogger("client")

	connection := Comm.NewConnection(DefaultConnectionURL, logger)
	connection.Listen()

	os.Exit(0)
}
