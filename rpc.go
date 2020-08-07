package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"shahrk/rpc/client"
	"shahrk/rpc/server"
)

var (
	port     = flag.Uint("port", 8080, "port to listen to")
	isServer = flag.Bool("server", false, "activate server mode")
	ip       = flag.String("ip", "127.0.0.1", "IP of server")
)

func handleSignal(s *server.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	sig := <-stop
	log.Println("exiting on receiving ", sig)
}

func startServer(s *server.Server) {
	defer stopServer(s)
	err := s.Start()
	if err != nil {
		log.Fatal("Error starting server:", err)
		return
	}
	handleSignal(s)
}

func stopServer(s *server.Server) {
	err := s.Stop()
	if err != nil {
		log.Fatal("Error while stopping server:", err)
	}
}

func testClient(c *client.Client) {
	log.Println("######################## Synchronuos Call ##############################")
	getArgs := &server.GetArgs{Key: "name"}
	putArgs := &server.PutArgs{Key: "name", Value: "Raj"}
	var putReply server.PutReply
	var getReply server.GetReply
	err := c.Client.Call("Vendy.Put", putArgs, &putReply)
	if err != nil {
		log.Fatal("Put Error:", err)
	} else if putReply.Err != "" {
		log.Fatal("Put Error:", putReply.Err)
	}
	log.Println("Put: " + putArgs.Key + "->" + putArgs.Value)
	log.Println("####################### Aysnchronous Call ##############################")
	getCall := c.Client.Go("Vendy.Get", getArgs, &getReply, nil)
	log.Println("printing asynchronously with rpc")
	getCall = <-getCall.Done
	if getCall.Error != nil {
		log.Fatal("Get Error:", getCall.Error)
	} else if getReply.Err != "" {
		log.Fatal("Get Error:", getReply.Err)
	} else {
		log.Println("Get: " + getArgs.Key + "->" + getReply.Value)
	}
}

func main() {
	flag.Parse()
	s := server.Server{Port: *port}
	c := client.Client{IP: *ip, Port: *port}

	if *isServer {
		startServer(&s)
		return
	}

	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting client to server:", err)
		return
	}
	testClient(&c)
}
