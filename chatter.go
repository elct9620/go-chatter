package main

import (
  "menteslibres.net/gosexy/redis"
  "code.google.com/p/go.net/websocket"
  "github.com/codegangsta/martini"
  "time"
  "fmt"
  "strings"
  "chatter/chat"
)

var (
  db *redis.Client
  clients map[string]*chatter.Client
)

func BroadcastMessage(sender *websocket.Conn, packet *chatter.Packet) {
  for _, item := range clients {
    websocket.Message.Send(item.Socket, packet.Pack())
  }
}

func ResponseMessage(sender *websocket.Conn, packet *chatter.Packet) {
  websocket.Message.Send(sender, packet.Pack())
}

func ParseCommand(rawMessage *string, client *chatter.Client) {
  if strings.Index(*rawMessage, "/") == 0 {
    commands := strings.SplitN(*rawMessage, " ", 2)

    var response string
    var broadcast string
    switch commands[0] {
    default:
      response = "No command found!"
    case "/nickname":
      broadcast = client.Nickname + " change nickname to " + commands[1]
      client.Nickname = commands[1]
    }
    if len(response) > 0 {
      ResponseMessage(client.Socket, &chatter.Packet{Type: "system", Data: response})
    }
    if len(broadcast) > 0 {
      BroadcastMessage(client.Socket, &chatter.Packet{Type: "system", Data: broadcast})
    }
    return
  }

  message := chatter.NewMessage(client.Nickname, *rawMessage)
  BroadcastMessage(client.Socket, &chatter.Packet{Type: "message", Data: message})
}

func wsHandler(ws *websocket.Conn) {
  var err error
  var clientMessage string

  defer ws.Close()

  uid := fmt.Sprintf("%d-%d", len(clients), time.Now().Unix())
  client := &chatter.Client{Nickname: "Guest" + uid, Socket: ws}

  clients[uid] = client
  defer func() {
    delete(clients, uid)
    client = nil
    fmt.Printf("Current online clients: %d\n", len(clients))
  }()

  for {
    if err = websocket.Message.Receive(ws, &clientMessage); err != nil {
      return
    }
    ParseCommand(&clientMessage, clients[uid])
  }

}

func main() {
  clients = make(map[string]*chatter.Client)

  m := martini.Classic()
  m.Use(martini.Static("public"))
  m.Get("/chatroom", websocket.Handler(wsHandler).ServeHTTP)
  m.Run()
}

func checkError(err error) {
  if err != nil {
    panic(err)
  }
}
