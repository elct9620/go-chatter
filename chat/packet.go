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
