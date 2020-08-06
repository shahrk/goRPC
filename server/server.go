package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type GetArgs struct {
	Key string
}

type GetReply struct {
	Err   string
	Value string
}

type PutArgs struct {
	Key, Value string
}

type PutReply struct {
	Err string
}

// Vendy is a thread-safe key value store
type Vendy struct {
	store map[string]string
	mux   sync.Mutex
}

func (v *Vendy) Get(args *GetArgs, reply *GetReply) (err error) {
	v.mux.Lock()
	defer v.mux.Unlock()
	if val, exists := v.store[args.Key]; !exists {
		log.Println("Unfulfilled GET Request")
		reply.Err = "Requested key does not exist"
	} else {
		log.Println("Fulfilled GET Request on Vendy")
		reply.Value = val
	}
	return
}

func (v *Vendy) Put(args *PutArgs, reply *PutReply) (err error) {
	v.mux.Lock()
	defer v.mux.Unlock()
	v.store[args.Key] = args.Value
	log.Println("Fulfilled PUT Request on Vendy")
	return
}

type Server struct {
	Port     uint
	listener net.Listener
}

func (server *Server) Stop() (err error) {
	if server.listener != nil {
		log.Println("GOODBYE")
		server.listener.Close()
	}
	return
}

func (server *Server) Start() (err error) {
	vendy := new(Vendy)
	vendy.store = make(map[string]string)
	err = rpc.Register(vendy)
	rpc.HandleHTTP()
	server.listener, err = net.Listen("tcp", fmt.Sprint(":", server.Port))

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
