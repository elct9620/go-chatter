package chatter

import (
  "code.google.com/p/go.net/websocket"
)

type Client struct {
  Nickname string
  Socket *websocket.Conn
}

