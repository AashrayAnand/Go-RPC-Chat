package tests

import (
  "../shared"
  "net/rpc"
  "testing"
  "time"
)

// given format parameter, returns stringified time
var t_ = time.Now().Format
var k_ = time.Kitchen

const (
  NEWUSER = "Server.NewUser" // RPC string for registering user
  GETMSG = "Server.GetMessages" // RPC string for getting messages
  SEND = "Server.Send" // RPC string for sending message
  EXIT = "/exit" // message user must send to exit chat
  HELP = "/help" // help function, user can send this to list chat commands
  WHITESPACE = ""
  local = "127.0.0.1:8080" // local host connection string
)

func TestRegisterUser(t *testing.T){

  // make RPC call directly, bypasses I/O stage of register
  conn, err := rpc.Dial("tcp", local) // register RPC client
  if err != nil {
    t.Fatalf("error connecting to server %s", err)
  }
  args := &shared.NewUserArgs{Name: "TEST_NAME"}
  resp := &shared.NewUserResp{}
  // RPC call to register new user
  err = conn.Call(NEWUSER, args, resp)
  if err != nil {
    t.Fatalf("error registering user")
  }

  // check if user was successfully registered
  args2 := &shared.GetMessagesArgs{User: "TEST_NAME"}
  resp2 := &shared.GetMessagesResp{}
  err = conn.Call(GETMSG, args2, resp2)
  if err == nil {
    t.Fatalf("user was not created")
  }
}

func TestRegisterDuplicate(t *testing.T){

    // make RPC call directly, bypasses I/O stage of register
    conn, err := rpc.Dial("tcp", local) // register RPC client
    if err != nil {
      t.Fatalf("error connecting to server %s", err)
    }
    args := &shared.NewUserArgs{Name: "TEST_NAME"}
    resp := &shared.NewUserResp{}
    // RPC call to register new user
    err = conn.Call(NEWUSER, args, resp)
    if err != nil {
      t.Fatalf("error registering user")
    }
    // duplicate register, should send response code of -1
    err = conn.Call(NEWUSER, args, resp)
    if resp.Code != -1 {
      t.Fatalf("error, duplicate registration allowed")
    }
}

func TestSendGroupMessage(t *testing.T){
  conn, _ := rpc.Dial("tcp", local) // register RPC client
  conn2, _ := rpc.Dial("tcp", local) // register RPC client
  args := &shared.SendMessageArgs{Message: &shared.Msg{Sender: "foo", Receiver:"", Message:"bar", Time: t_(k_)}}
  resp := &shared.SendMessageResp{}
  err := client.Conn.Call(SEND, args, resp)
  // error check for robustness, server does not return
  // non-nil error for any case, but can leave this here
  // in case server send RPC stub is updated
  if err != nil {
    t.Fatalf(err)
  }

  args2 := &shared.GetMessagesArgs{User: "TEST_NAME"}
  resp2 := &shared.GetMessagesResp{}
  err = conn.Call(GETMSG, args2, resp2)
  if err == nil {
    t.Fatalf("user was not created")
  }
}

/*func TestSendDirectMessage(t *testing.T){
  args := &shared.SendMessageArgs{Message: message}
  resp := &shared.SendMessageResp{}
  err := client.Conn.Call(SEND, args, resp)
  // error check for robustness, server does not return
  // non-nil error for any case, but can leave this here
  // in case server send RPC stub is updated
  if err != nil {
    log.Fatal(err)
  }
  // attempted to DM user not currently in chat room
  if resp.Code == -1 {
    fmt.Printf("the user %s does not exist\n", message.Receiver)
  }
}*/
