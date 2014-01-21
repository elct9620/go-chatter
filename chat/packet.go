package chatter

import (
  "github.com/vmihailenco/msgpack"
)

type Packet struct {
  Type string
  Data interface{}
}

func (p Packet) Pack() []byte {
  packet, err := msgpack.Marshal(p)
  if err != nil {
    return nil
  }
  return packet
}

/**
 * Packet Maker
 */

func NewSystemPacket(message string) *Packet {
  return &Packet{Type: "system", Data: message}
}

func NewMessagePacket(nickname string, content string) *Packet {
  return &Packet{Type: "message", Data: NewMessage(nickname, content)}
}

func NewOnlineCountPacket(count int) *Packet {
  return &Packet{Type: "online-count", Data: count}
}


