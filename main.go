package main

import (
	"fmt"
	"flag" // used to get command line arguments
	"log"
	"time"
	"./client"
	"./server"
)

// given format parameter, returns stringified time
var t_ = time.Now().Format
var k_ = time.Kitchen

func main() {
	var proc string
	flag.StringVar(&proc, "proc", "", "") // defines -proc flag, default value of server
	flag.Parse() // parse command line arguments

	if(proc != "server" && proc != "client"){
		log.Fatal("error, -proc argument must be server or client")
	} else {
		fmt.Println("initiating", proc,"process...")
		if (proc == "server"){
			srv := new(server.Server) // create server structure
			srv.Serve() // register RPC structs/methods and listen to local:8080
			defer srv.Terminate()
		} else {
			cli := new(client.Client) // create client structure
			cli.Create() // initialize client structure
			cli.Handle() // register user and handle sending/receiving messages
		}
	}
}
