package main

import (
  "menteslibres.net/gosexy/redis"
  "code.google.com/p/go.net/websocket"
  "github.com/codegangsta/martini"
  "time"
  "fmt"
  "strings"
  "chatter/chat"
  "chatter/helper"
)

var (
  db *redis.Client
  clients map[string]*chatter.Client
  logSize int64 = 50
)

func Broadcast(sender *websocket.Conn, packet *chatter.Packet) {
  for _, item := range clients {
    websocket.Message.Send(item.Socket, packet.Pack())
  }
  // Store chat log
  switch packet.Type {
  case "message", "system":
    db.RPush("chat:logs", packet)
    if logLen, _ := db.LLen("chat:logs"); logLen > logSize {
      log, _ := db.LPop("chat:logs")
      fmt.Println("Remove log: " + log)
    }
  }
}

func Response(sender *websocket.Conn, packet *chatter.Packet) {
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
      Response(client.Socket, chatter.NewSystemPacket(response))
    }
    if len(broadcast) > 0 {
      Broadcast(client.Socket, chatter.NewSystemPacket(broadcast))
    }
    return
  }

  Broadcast(client.Socket, chatter.NewMessagePacket(client.Nickname, *rawMessage))
}

func wsHandler(ws *websocket.Conn) {
  var err error
  var clientMessage string

  defer ws.Close()

  uid := fmt.Sprintf("%d-%d", len(clients), time.Now().Unix())
  client := &chatter.Client{Nickname: "Guest" + uid, Socket: ws}
  clients[uid] = client

  Broadcast(client.Socket, chatter.NewOnlineCountPacket(len(clients)))

  defer func() {
    delete(clients, uid)
    clientCount := len(clients)
    Broadcast(client.Socket, chatter.NewOnlineCountPacket(clientCount))
    client = nil
    fmt.Printf("Current online clients: %d\n", clientCount)
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

  db = redis.New()
  dbHost, dbPort, dbPassword := helper.GetRedisToGoEnv()
  db.Connect(dbHost, dbPort)
  if len(dbPassword) > 0 {
    db.Auth(dbPassword)
  }

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
