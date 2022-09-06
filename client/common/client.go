package common

import (
	"net"
	"time"
	"bufio"
	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Record used by the client to represent a person
type PersonRecord struct {
	FirstName string
	LastName string
	Document uint64
	Birthdate string
}

// Client Entity that encapsulates how
type Client struct {
	person PersonRecord
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, person PersonRecord) *Client {
	client := &Client{
		person: person,
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
func (c *Client) Start() {
	c.createClientSocket()
	
	// Send
	log.Infof("Sending person query")
	err := Send(c.conn, AskWinner, c.person)
	if err != nil{
		log.Errorf("[CLIENT %v] %v", c.config.ID, err)
		c.conn.Close()
		return
	}

	log.Infof("Awaiting server answer")
	reader := bufio.NewReader(c.conn)
	isWinner, err := RecvBool(reader)
	if err != nil{
		log.Errorf("[CLIENT %v] %v", c.config.ID, err)
		c.conn.Close()
		return
	}

	if isWinner{
		log.Infof("It's a winner")
	}else{
		log.Infof("It isn't a winner")
	}
	
	c.conn.Close()
}
