package core

import (
  "fmt"
)

// utility function, takes a slice length and some function,
// and checks each element of the slice with the specified function
// returning the index that passes the check or -1
func CheckSlice(length int, check func(i int) bool){
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

type NewUserArgs struct { // args struct for REQUESTING new user
      Name string // name of new user
}

// RPC #2: SEND MESSAGE

/*type SendMessageArgs struct { // args struct for REQUESTING message be sent
  message Msg // message being sent
}

type SendMessageResp struct {

}*/

// RPC #3: GET MESSAGES

/*type GetMessagesArgs struct { // args struct for to REQUESTING user's unseen messages
  User string // user to get messages for
  Index int // position of user in messsage queue, as of most recent check
}

type GetMessagesResp struct { // response struct for RECEIVING user's unseen messages
  Messages []Msg // messages unseen by user
  Index int // user position in message queue after getting lastest messages
}*/

// stringify a message struct
func (m *Msg) String() string {
  return fmt.Sprintf("[%s @ %s]: %s", m.Name, m.Time, m.Message)
}
