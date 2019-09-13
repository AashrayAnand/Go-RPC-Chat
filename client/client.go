package client

import (
  "../shared"
  "net/rpc"
  "time"
  "log"
  "fmt"
  "bufio"
  "os"
  "strings"
  "sync"
)

// given format parameter, returns stringified time
var t_ = time.Now().Format
var k_ = time.Kitchen

type Client struct {
        Name string // screen name
        Conn *rpc.Client // result of calling rpc.Dial
}

const (
  NEWUSER = "Server.NewUser" // RPC string for registering user
  GETMSG = "Server.GetMessages" // RPC string for getting messages
  SEND = "Server.Send" // RPC string for sending message
  EXIT = "/exit" // message user must send to exit chat
  HELP = "/help" // help function, user can send this to list chat commands
  WHITESPACE = ""
  host = "18.219.140.44:" // local host connection string
  port = "3000"
)

func (client *Client) Terminate() {
  if client.Conn != nil {
    client.Conn.Close()
  }
}

func (client *Client) Create() {
  if client.Conn == nil {
    var err error
    client.Conn, err = rpc.Dial("tcp", host + port) // register RPC client
    if err != nil {
      log.Fatal("client error on Dial")
    }
  }
}

// client-side method, executes RPC call to get unseen messsages, matches
// name of remote function for transparency
func (client *Client) GetMessages(DoneChan chan int) {
  go func () {
    <- DoneChan
    return
  }()
  for {
    args := &shared.GetMessagesArgs{User: client.Name}
    resp := &shared.GetMessagesResp{}
    err := client.Conn.Call(GETMSG, args, resp)
    if err != nil {
      log.Fatal(err)
    }

    // display unseen messages
    for _, message := range resp.Messages {
      fmt.Println(message)
    }
  }
}

func (client *Client) SendMessage(message *shared.Msg) {
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
}

// listener for client

// main loop for client, dispatch goroutines to check for and handle messages
func (client *Client) Handle() {
  // use channel to block Handle() after dispatching go routines to get
  // messages and listen to user input (value will be sent to DoneChan
  // when client.listen terminates)
  DoneChan := make(chan int)
  // terminate when Handle() exits
  defer client.Terminate()
  var wg sync.WaitGroup
  wg.Add(1)
  go client.register(&wg) // register user in chat room
  wg.Wait()

  go client.GetMessages(DoneChan) // get user's messages
  go client.listen(DoneChan) // listen for user input
  <- DoneChan
}

// asks for screen name continuously, until error, or name is registered (unique name)
func (client *Client) register(wg *sync.WaitGroup) {
  reader := bufio.NewReader(os.Stdin) // read from stdin
  for {
    fmt.Print("Enter a screen name: ")
    name, _ := reader.ReadString('\n')
    nameNoLine := name[:len(name) - 1]
    if strings.IndexByte(nameNoLine, ' ') != -1 {
      fmt.Println("cannot have white space in user name")
      continue
    }
    fmt.Printf("attempting to register user: %s", name)
    args := &shared.NewUserArgs{Name: nameNoLine}
    resp := &shared.NewUserResp{}
    // RPC call to register new user
    err := client.Conn.Call(NEWUSER, args, resp)
    // err is only non-nil if chat room is full, should exit
    if err != nil {
      log.Fatal(err)
    }
    // name is taken, failed to register user
    if resp.Code == -1 {
      fmt.Printf("Failed to register, name %s is taken\n", name[:len(name) - 1])
      continue
    // successfully registered user
    } else {
      fmt.Printf("Successfully registered user, welcome %s", name)
      client.Name = nameNoLine
      break
    }
  }
  wg.Done()
}

func (client *Client) listen(DoneChan chan int) {
  reader := bufio.NewReader(os.Stdin)
  for {
    // wait for message
    message, _ := reader.ReadString('\n')
    // remove newline from message
    messageNoLine := message[:len(message) - 1]
    switch messageNoLine {
    case WHITESPACE:
      continue
    // list commands supported by the chat
    case HELP:
      fmt.Println("CHAT COMMANDS")
      fmt.Println("exit the chat -> /exit")
      fmt.Println("send a DM to recipient -> @recipient <message>")
    // send
    case EXIT:
      DoneChan <- 1
      break
    default:
      // direct message
      if messageNoLine[0] == '@' {
        directMessage := messageNoLine[1:]
        // receiver and message should be split by white space
        // get the position of the split
        endOfReceiver := strings.IndexByte(directMessage, ' ')
        if endOfReceiver == -1 || endOfReceiver == len(directMessage) - 1 {
          fmt.Println("please enter non-empty message")
          continue
        }
        // get receiver (characters before whitespace)
        receiver := directMessage[:endOfReceiver]
        // get message (characters after whitespace)
        message := directMessage[endOfReceiver + 1:]
        // function to check if message is just white space
        checkIsWhitespace := func(s string) bool {
          for i := 0; i < len(s); i++ {
            if s[i] != ' ' {
              return false
            }
          }
          return true
        }
        // send message, if not just white space
        if !checkIsWhitespace(message) {
          MessageStruct := &shared.Msg{Sender: client.Name, Receiver : receiver, Message: message, Time: t_(k_)}
          client.SendMessage(MessageStruct)
        } else {
          fmt.Println("message must contain non-whitespace characters")
        }
      } else {
        // send basic message to group
        MessageStruct :=  &shared.Msg{Sender: client.Name, Receiver: "", Message: messageNoLine, Time: t_(k_)}
        client.SendMessage(MessageStruct)
      }
    }
  }
}
