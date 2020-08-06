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
	defer s.Stop()
	s.Start()
	handleSignal(s)
}

func testClient(c *client.Client) (err error) {
	log.Println("######################## Synchronuos Call ##############################")
	args := &server.Args{30, 10}
	var reply int
	err = c.Client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith multiply error:", err)
	}
	log.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)
	log.Println("####################### Aysnchronous Call ##############################")
	quotient := new(server.Quotient)
	divCall := c.Client.Go("Arith.Divide", args, quotient, nil)
	log.Println("printing asynchronously with rpc")
	divCall = <-divCall.Done
	if divCall.Error != nil {
		log.Fatal("arith div error:", divCall.Error)
	} else {
		log.Printf("Arith: %d/%d=%dx + %d", args.A, args.B, quotient.Quo, quotient.Rem)
	}
	return
}

func main() {
	flag.Parse()
	s := server.Server{Port: *port}
	c := client.Client{IP: *ip, Port: *port}

	if *isServer {
		startServer(&s)
		return
	}

	c.Connect()
	testClient(&c)
}
