package main

import (
	"fmt"
	"flag" // used to get command line arguments
	"log"
	"time"
	"./core"
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
			cli := new(client.Client) // create client structureßß ß
			cli.Connect()
			//cli.Message(&core.Msg{Name: "A", Message: "Hello", Time: t_(k_)})
			cli.Handle()
			defer cli.Terminate() // close connection on termination
		}
	}
}
