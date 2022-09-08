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

// Client Entity that encapsulates how
type Client struct {
	persons []PersonRecord
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, personRecords []PersonRecord) *Client {
	client := &Client{
		persons: personRecords,
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

func sendAll(socket net.Conn, msg []byte) error{
	nwritten_acum := len(msg)
	for ;nwritten_acum > 0; {
		nwritten, err := socket.Write(msg)
		if err != nil {
			return err
		}
		nwritten_acum -= nwritten
	}
	return nil
}

func (c *Client) askWinner(){
log.Infof("Sending person query")
	msg := SerializeUint32(AskWinner)
	err := sendAll(c.conn, msg)
	if err != nil{
		log.Errorf("[CLIENT %v] %v", c.config.ID, err)
		c.conn.Close()
		return
	}

	msg = SerializePersonRecordArray(c.persons)
	err = sendAll(c.conn, msg)
	if err != nil{
		log.Errorf("[CLIENT %v] %v", c.config.ID, err)
		c.conn.Close()
		return
	}

	log.Infof("Awaiting server answer")
	winners := 0
	reader := bufio.NewReader(c.conn)
	for i := 0; i < len(c.persons); i++ {
		isWinner, err := DeserializeBool(reader)
		if err != nil{
			log.Errorf("[CLIENT %v] %v", c.config.ID, err)
			c.conn.Close()
			return
		}

		if isWinner{
			winners++
		}
	}

	log.Infof("Records sent: %v, Winners: %v %%", len(c.persons), (winners*100)/len(c.persons))
}

func (c *Client) askAmount(){
	reader := bufio.NewReader(c.conn)
	
	for ;true;{
		log.Infof("Sending amount query")
		msg := SerializeUint32(AskAmount)
		err := sendAll(c.conn, msg)
		if err != nil{
			log.Errorf("[CLIENT %v] %v", c.config.ID, err)
			c.conn.Close()
			return
		}

		log.Infof("Awaiting server answer")
		total, err := DeserializeUint32(reader)
		if err != nil{
			log.Errorf("[CLIENT %v] %v", c.config.ID, err)
			c.conn.Close()
			return
		}
		partial, err := DeserializeBool(reader)
		if err != nil{
			log.Errorf("[CLIENT %v] %v", c.config.ID, err)
			c.conn.Close()
			return
		}
		log.Infof("Total winners: %v, Partial: %v", total, partial)
		if !partial {
			break
		}
		time.Sleep(8 * time.Second)
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) Start() {
	c.createClientSocket()
	c.askWinner()
	c.askAmount()
	c.conn.Close()
}
