package server
import (
  "net/rpc"
  "../shared"
  "net"
  "strconv"
//  "net/http"
  "log"
  "fmt"
  "errors"
  "sync"
)

// wrapped for message list dictionary, allows for atomicity
// when reading/writing messsages
type LockedMessages struct {
  Map map[string][]*shared.Msg
  sync.Mutex
}

type Server struct {
  Port int // port to connect to
  Users []string // users in the chat room
  Messages LockedMessages // map of each user to their message list
  listener net.Listener // listener for RPC calls
}

const MAX_USERS = 100 // at most 100 people in the chat room

// exported Server methods should be able to
// 1. register a new user
// 2. broadcast a message to the chat
// 3. broadcast a message to a specific user
// 4. terminate the chat

// utility method, calls generic shared slice indexing function, for
// this method, simply returns index of a name in server's user list
// or -1 if user does not exist
func (server *Server) checkName(name string) int {
  // function passed to CheckSplice to check if name is duplicate
  predicate := func (i int) bool {
    return server.Users[i] == name
  }

  return shared.CheckSlice(len(server.Users), predicate)
}

// exported method, will be executed through RPC by clients attempting to join chat room
func (server *Server) NewUser(request *shared.NewUserArgs, response *shared.NewUserResp) error {
  if len(server.Users) == MAX_USERS {
    return errors.New("Chat room is full. Try again later.")
  }

  if server.checkName(request.Name) != -1 {
    fmt.Println("blocked attempt to register user with an existing name", request.Name)
    response.Code = -1
  } else {
    server.Users = append(server.Users, request.Name) // add user to list of users
    fmt.Println("registered new user", request.Name)
    response.Code = 0
  }
  return nil
}

// exported method, will be executed through RPC by clients attempting to get messages
func (server *Server) GetMessages(request *shared.GetMessagesArgs, response *shared.GetMessagesResp) error {
  // lock message list until we have copied messages
  server.Messages.Lock()
  defer server.Messages.Unlock()
  if messages, ok := server.Messages.Map[request.User]; ok {
    response.Messages = messages
  } else {
    return errors.New("user does not exist")
  }
  server.Messages.Map[request.User] = nil
  return nil
}

// exported method, will be executed through RPC by clients attempting to send a message
func (server *Server) Send(request *shared.SendMessageArgs, response *shared.SendMessageResp) error {
  // message is a group message, send to all message lists
  fmt.Println("received message", request.Message)
  server.Messages.Lock()
  if request.Message.Receiver == "" {
    for user, _ := range server.Messages.Map {
      server.Messages.Map[user] = append(server.Messages.Map[user], request.Message)
    }
    response.Code = 0
  // direct message, check if user exists, if so, send to user's message list
  } else {
    receiver := request.Message.Receiver
    // user exists, add to their message list
    if server.checkName(receiver) != -1 {
      server.Messages.Map[receiver] = append(server.Messages.Map[receiver], request.Message)
      response.Code = 0
    } else {
      fmt.Println("attempted to send message to non-existent user", receiver)
      response.Code = -1
    }
  }
  server.Messages.Unlock()
  return nil
}

func (server *Server) Terminate() {
  if server.listener != nil {
    server.listener.Close()
  }
}

func (server *Server) initialize() {
  server.Port = 8080
  server.Users = make([]string, 1, MAX_USERS) // initialize to hold up to max users
  server.Messages.Map = make(map[string][]*shared.Msg, 1) // will be resized with append
}

// create a chat server
func (server *Server) Serve() {
  server.initialize() // initialize server structure fields
  rpc.Register(server) // register server methods which satisfy RPC constraints

  var err error
  server.listener, err = net.Listen("tcp", ":"+strconv.Itoa(server.Port)) // return net listener on port 8080

  if err != nil { // error checking
    log.Fatal("server listen error:", err)
  }

  rpc.Accept(server.listener)
}
