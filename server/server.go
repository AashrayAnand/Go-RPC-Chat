package server
import (
  "net/rpc"
  "../core"
  "net"
  "net/http"
  "log"
  "fmt"
)

type Server struct {
  Users []string // users in the chat room
  UserMessages map[string][]string // for each user, maintain a queue of messages
  listener net.Listener
}

// exported Server methods should be able to
// 1. register a new user
// 2. broadcast a message to the chat
// 3. broadcast a message to a specific user
// 4. terminate the chat

func (server *Server) GetMessages(request *core.GetMessagesArgs, response *core.MsgListResp) error {
  // check message queue exists for user,
  // check length of message queue is > curr index of client
  // return slice of message queue from curr index of client onwards
  if messages, ok := server.UserMessages[request.User]; ok {
    // store messages since last check, update current index
    response.Messages = messages[request.Index:]
    if len(messages) == 0 {
      response.Index = 0
    } else {
      response.Index = len(messages) - 1
    }
    return nil
  }
  return error.New("User %s does not exist!", request.User)
}

// server echoes message
func (server *Server) Echo(request *core.Msg, response *core.MsgResp) error {
  // create response message for message
  response.Message = fmt.Sprintf("[%s @ %s]: %s", request.Name, request.Time, request.Message)
  return nil
}

func (server *Server) Terminate() {
  if server.listener != nil {
    server.listener.Close()
  }
}

// create a chat server
func (server *Server) Serve() {
  rpc.Register(server) // register server methods which satisfy RPC constraints
  rpc.HandleHTTP() // indicate that RPC server receives HTTP requests

  listener, err := net.Listen("tcp", ":8080") // return net listener on port 8080

  if err != nil { // error checking
    log.Fatal("server listen error:", err)
  }

  server.listener = listener // bind net listener to server struct


  http.Serve(server.listener, nil) // listen on net listener, and dispatch
                                   // go routines to service requests
}
