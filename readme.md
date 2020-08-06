# A simple RPC client/server

This is the code for a simple RPC client/server I created to try out the RPC functionalities GO provides.


## What does it do?

The _server_ has a thread safe key value store called vendy which implements two methods - Get & Put, it has a start method which registers the args struct with rpc and starts listening on the provided port. It also has a stop method which closes the listener.

The _client_ has a connect method which dials to the server's ip/port.

The program when run (default mode) runs the client's connect method and tries to make RPCs to the Put & Get methods. When run in server mode it starts the server and listens for interrupt on getting which is calls the server's stop method


## How do I set it up?

Once you have pulled the code run the following to install the module

```
go install
```

This should create the executable rpc in $GOPATH/bin/. If you have the bin added to your PATH variable you should now directly be able to run the program.

The program runs in client mode by default. To run the program in server mode just pass the -server flag.

```
rpc -server
```

Checkout the help to get familiar with the flags. Run accordingly. 

```
rpc --help
```