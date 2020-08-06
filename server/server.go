package server

import (
	"errors"
	"net/rpc"
	"net"
	"log"
	"net/http"
	"fmt"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

type Server struct {
	Port uint
	listener net.Listener
}

func (server *Server) Stop() (err error) {
	if server.listener != nil {
		server.listener.Close()
	}
	return
}

func (server *Server) Start() (err error) {
	if server.Port < 0 {
		err = errors.New("Port is missing!")
		return
	}
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	server.listener, err = net.Listen("tcp", fmt.Sprint(":",server.Port))

	if err != nil {
		log.Fatal("Listen error:", err)
		return
	} else {
		log.Println("Server listening on port ", server.Port)
	}
	go func() {
		if err := http.Serve(server.listener, nil); err != nil {
			log.Fatal("Serve error:", err)
		}
	}()
	return
}
