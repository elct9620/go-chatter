package chatter

import (
  "github.com/vmihailenco/msgpack"
  "fmt"
  "strings"
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

func (p Packet) String() string {
  switch p.Data.(type) {
    case string, fmt.Stringer:
      return fmt.Sprintf("%s::%s", p.Type, p.Data)
  }
  return ""
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

func NewLogsPacket(packets []*Packet) *Packet {
  return &Packet{Type: "logs", Data: packets}
}

/**
 * Upacket
 */

func StringToPacket(rawData string) *Packet {
  rawPacket := strings.SplitN(rawData, "::", 2)
  if len(rawPacket) != 2 {
    return nil
  }
  switch rawPacket[0] { // Check packet type
    case "system":
      return &Packet{Type: rawPacket[0], Data: rawPacket[1]}
    case "message":
      return &Packet{Type: rawPacket[0], Data: StringToMessage(rawPacket[1])}
  }
  return nil
}
