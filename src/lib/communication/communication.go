// Communicate via Mangos sockets.
//
// Exports
// -------
// type Connection struct
//
// NewConnection(
// 	url string,
// 	mode string,
// 	logger *log.logger,
// 	debug bool,
// 	[msgHandler func([]byte)]
// ) *Connection: Create a new connection.
//
// Connection.Run(): Initialize a connection.
// Connection.Send([]byte): Send a byte string via the connection.
// Connection.Recv() []byte: Receive a byte string via the conncetion.

package communication

import (
	"fmt"
	"log"
	"time"

	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/pair"
	"nanomsg.org/go-mangos/transport/tcp"
)

type Connection struct {
	url         string
	socket      mangos.Socket
	logger      *log.Logger
	connRoutine func()
	msgHandler  func([]byte)
	debug       bool
}

// logError logs an error variable via the connections logger.
func (conn Connection) logError(err error) {
	conn.logger.Println("error: ", fmt.Sprintf("can't listen on pair socket: %s", err.Error()))
}

// NewConnection allocates a new Connection, sets up the logger, connection url, debug mode, creates a new socket,
// select the correct routine to run, and returns the pointer to the Connection.
func NewConnection(url string, mode string, logger *log.Logger, debug bool, optional ...interface{}) *Connection {
	conn := new(Connection)

	conn.logger = logger
	conn.url = url
	conn.debug = debug

	if len(optional) > 0 {
		conn.msgHandler = optional[0].(func([]byte))
	}

	var err error
	if conn.socket, err = pair.NewSocket(); err != nil {
		conn.logError(err)
	}

	conn.socket.AddTransport(tcp.NewTransport())

	switch mode {
	case "PAIR.SERVER":
		conn.connRoutine = conn.server
	case "PAIR.CLIENT":
		conn.connRoutine = conn.client
	case "PAIR":
		conn.connRoutine = conn.client
	}

	return conn
}

// server listens for messages on the mangos socket bound to the connection and executes
// the message handler when a message is received.
func (conn Connection) server() {
	var err error
	var recvd []byte
	if err = conn.socket.Listen(conn.url); err != nil {
		conn.logger.Println("error: ", fmt.Sprintf("can't listen on pair socket: %s", err.Error()))
	}

	for {
		conn.socket.SetOption(mangos.OptionRecvDeadline, 100*time.Millisecond)
		if recvd, err = conn.Recv(); err == nil {
			conn.msgHandler(recvd)
		}

	}
}

// client dials the TCP url.
func (conn Connection) client() {
	var err error
	if err = conn.socket.Dial(conn.url); err != nil {
		conn.logger.Println("error: ", fmt.Sprintf("can't dial on pair socket: %s", err.Error()))
	}
}

// Run is the exported routine that executes the correct function specfic
// to the the socket type as selected by NewConnection.
func (conn Connection) Run() {
	conn.connRoutine()
}

// Send transmits a byte string using the mangos socket.
func (conn Connection) Send(msg []byte) {
	var err error
	if err = conn.socket.Send(msg); err != nil {
		conn.logger.Println("error: ", fmt.Sprintf("failed sending: %s", err.Error()))
	}
	if conn.debug {
		conn.logger.Println("send: ", string(msg))
	}
}

// Recv receives and returns the byte string received by the connection.
func (conn Connection) Recv() ([]byte, error) {
	var msg []byte
	var err error
	if msg, err = conn.socket.Recv(); err == nil {
		if conn.debug {
			conn.logger.Println("recvd: ", string(msg))
		}
		return msg, err
	}

	return msg, err
}
