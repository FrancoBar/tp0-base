package common

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	signals chan os.Signal
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, signals chan os.Signal) *Client {
	client := &Client{
		signals: signals,
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"[CLIENT %v] Could not connect to server. Error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	//Send an autoincremental msgID to identify every message sent
	msgID := 1

loop:
	// Send messages if the loopLapse threshold has been not surpassed
	for timeout := time.After(c.config.LoopLapse); ; {

		// Create the connection the server in every loop iteration. 
		c.createClientSocket()
		select {
		case <-timeout:
			break loop
		case <- c.signals:
			log.Infof("[CLIENT %v] SIGTERM or SIGINT received.", c.config.ID)
			break loop
		// Wait a time between sending one message and the next one
		case <-time.After(c.config.LoopPeriod):
		}

		// Send
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v sent\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msgID++

		if err != nil {
			log.Errorf(
				"[CLIENT %v] Error reading from socket. %v.",
				c.config.ID,
				err,
			)
			c.conn.Close()
			return
		}
		log.Infof("[CLIENT %v] Message from server: %v", c.config.ID, msg)

		// Recreate connection to the server
		c.conn.Close()
	}

	log.Infof("[CLIENT %v] Closing connection", c.config.ID)
	c.conn.Close()
}
