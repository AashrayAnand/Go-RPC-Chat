package client

import (
  "../core"
  "net/rpc"
  "time"
  "log"
  "fmt"
)

// given format parameter, returns stringified time
var t_ = time.Now().Format
var k_ = time.Kitchen

type Client struct {
        Name string // screen name
        Index int // client's current positon in message queue
        Conn *rpc.Client // result of calling rpc.Dial
}

const (
  GETMSG = "Server.GetMessages"
  SEND = "Server.Send"
  local = "127.0.0.1:8080" // local host connection string
)

// use RPC to send message
func (client *Client) Message(message *core.Msg) {
  var response core.MsgResp

  err := client.Conn.Call(ECHO, message, &response) // uses RPC to call Handler.Send with given args
  if err != nil {
    log.Fatal("error on Message(...):", err)
  }

  fmt.Printf("%s\n", response.Message)
}

func (client *Client) Terminate() {
  if client.Conn != nil {
    client.Conn.Close()
  }
}

func (client *Client) Create() {
  if client.Conn == nil {
    var err error
    client.Conn, err = rpc.DialHTTP("tcp", local) // local host HTTP RPC server connection
    if err != nil {
      log.Fatal("client error on Dial")
    }
  }
}

func (client *Client) HandleMessages(ch chan int) {
  // TODO: loops infinitely, handling receiving and sending messages
  // by client
  ch <- 1
}

// main loop for client, dispatch goroutines to check for and handle messages
func (client *Client) Handle() {
  go func() {
    for {
      args := &core.GetMessagesArgs{User: client.Name, Index: client.Index}
      response := &core.MsgListResp{} // response (list of messages)
      err := client.Conn.Call(GETMSG, args, response)
      if err != nil {
        log.Fatal("error checking for messages", err)
      }
      for _ , msg := range response.Messages {
        fmt.Println(msg)
      }
    }
  }()

  ch := make(chan int)
  go client.HandleMessages(ch)
  <- ch // blocks on receiving a value from channel, only on
        // termination of HandleMessages
}
