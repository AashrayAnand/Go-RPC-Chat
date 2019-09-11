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
  Messages []core.Msg // list of message
  listener net.Listener // HTTP listener
}

// exported Server methods should be able to
// 1. register a new user
// 2. broadcast a message to the chat
// 3. broadcast a message to a specific user
// 4. terminate the chat

func (server *Server) NewUser(request *core.NewUserArgs, response *int) error {
  if server.Users == nil {
    server.Users = make([]string, 1) // create slice to store users
  }
  // utility function checks if user name already exists
  if core.CheckSlice(len(server.Users,
    func (i int) bool { return server.Users[i] == request.Name})) != -1 {
    *response = -1
  } else {
    server.Users = append(server.Users, request.Name) // add user to list of users
    *response = 0
  }
  return nil // no errors (duplicate name just returns appropriate response)
}

/*func (server *Server) GetMessages(request *core.GetMessagesArgs, response *core.GetMessagesResp) error {
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
}*/

func (server *Server) Terminate() {
  if server.listener != nil {
    server.listener.Close()
  }
}

// create a chat server
func (server *Server) Serve() {
  rpc.Register(server) // register server methods which satisfy RPC constraints
  rpc.HandleHTTP() // indicate that RPC server receives HTTP requests

  var err error
  server.listener, err = net.Listen("tcp", ":8080") // return net listener on port 8080

  if err != nil { // error checking
    log.Fatal("server listen error:", err)
  }

  http.Serve(server.listener, nil) // listen on net listener, and dispatch
                                   // go routines to service requests
}
