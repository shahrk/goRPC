package client

import (
	"log"
	"net/rpc"
	"fmt"
)

type Client struct {
	IP string
	Port uint
	Client *rpc.Client
}

func (c *Client) Connect() (err error) {
	c.Client, err = rpc.DialHTTP("tcp", fmt.Sprint(c.IP,":",c.Port))
	if (err != nil) {
		log.Fatal("Dial error:", err)
	} else {
		log.Println("Client connected to " + fmt.Sprint(c.IP,":",c.Port))
	}
	return
}
