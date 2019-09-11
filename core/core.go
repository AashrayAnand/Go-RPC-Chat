package core

import (
  "fmt"
)

// struct for representing a chat message
type Msg struct {
      // sender, message contents, time of message
      Name, Message, Time string
}

type MsgResp struct {
      Message string
}

type GetMessagesArgs struct { // passed to server, to get user's unseen messages
  User string // user to get messages for
  Index int // position of user in messsage queue (used to deem messages as seen or unseen)
}

type MsgListResp struct { // response when client checks for new messages
  Messages Msg[]
  Index int // user position in message queue after getting lastest messages
}

// stringify a message struct
func (m *Msg) String() string {
  return fmt.Sprintf("[%s @ %s]: %s", m.Name, m.Time, m.Message)
}
