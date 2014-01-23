package chatter

import (
  "strings"
)

type Message struct {
  Name string
  Content string
}

func (m Message) String() string {
  return m.Name + ":" + m.Content
}

func NewMessage(name string, content string) *Message {
  return &Message{Name: name, Content: content}
}

func StringToMessage(rawData string) *Message {
  rawMessage := strings.SplitN(rawData, ":", 2)
  if len(rawMessage) != 2 {
    return nil
  }
  return NewMessage(rawMessage[0], rawMessage[1])
}
