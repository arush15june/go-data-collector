package main

//                 Server
// Establish listener on DefaultConnectionUrl,
// Handle received messages via msgHandler,
// msgHandler: log recvd messages

import (
	"log"
	"os"

	Comm "../lib/communication"
	Logger "../lib/logger"
)

const DefaultConnectionUrl = "tcp://0.0.0.0:1337"

var (
	logger     *log.Logger
	connection *Comm.Connection
)

func msgHandler(msg []byte) {
	logger.Println("recvd: ", string(msg))
}

func main() {
	// Initilaize logger, listener
	logger = Logger.NewLogger("server")

	logger.Println("Starting Server")
	connection := Comm.NewConnection(DefaultConnectionUrl, "PAIR.SERVER", logger, false, msgHandler)
	connection.Run()

	os.Exit(0)
}
