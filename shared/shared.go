package shared

import (
  "fmt"
)

// utility function, takes a slice length and some function,
// and checks each element of the slice with the specified function
// returning the index that passes the check or -1
func CheckSlice(length int, check func(i int) bool) int {
  for i := 0; i < length; i++ {
    if check(i) {
      return i
    }
  }
  return -1 // no element passed check
}

type Msg struct { // struct representing a chat message
      // users messages is sent to, message contents, time of message
      Sender, Receiver, Message, Time string
}

// RPC #1: NEW USER

// args struct for REQUESTING new user
type NewUserArgs struct {
      Name string // name of new user
}

// response struct for RECEIVING result of requesting
// new user (integer wrapper, for clarity only)
type NewUserResp struct {
      Code int
}

// args struct for REQUESTING new user
type RemoveUserArgs struct {
      Name string // name of new user
}

// response struct for RECEIVING result of requesting
// new user (integer wrapper, for clarity only)
type RemoveUserResp struct {
      Code int
}

// RPC #2: SEND MESSAGE

// args struct for requesting message be sent
type SendMessageArgs struct {
      Message *Msg // message being sent
}

// response struct for RECEIVING result of sending
// message (integer wrapper, for clarity only)
type SendMessageResp struct {
      Code int
}

// RPC #3: GET MESSAGES

// args struct for to REQUESTING user's unseen messages
type GetMessagesArgs struct {
  User string // user to get messages for
}

// response struct for RECEIVING user's unseen messages
type GetMessagesResp struct {
  Messages []*Msg // messages unseen by user
}

// stringify a message struct
func (m *Msg) String() string {
  return fmt.Sprintf("[%s @ %s]: %s", m.Sender, m.Time, m.Message)
}
