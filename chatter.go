package main

import (
  "menteslibres.net/gosexy/redis"
  "code.google.com/p/go.net/websocket"
  "github.com/codegangsta/martini"
  "container/list"
  "time"
  "fmt"
  "strings"
  "chatter/chat"
)

var (
  db *redis.Client
  clients *list.List
  nickname map[string]string
)

func BroadcastMessage(sender *websocket.Conn, message *string) {
  for item := clients.Front(); item != nil; item = item.Next() {
    socket, ok := item.Value.(*websocket.Conn)
    if !ok {
      panic("it isn't *websocket.Conn")
    }

    websocket.Message.Send(socket, *message)
  }
}

func ResponseMessage(sender *websocket.Conn, message *string) {
  websocket.Message.Send(sender, *message)
}

func ParseCommand(rawMessage *string, ws *websocket.Conn, uid *string) {
  if strings.Index(*rawMessage, "/") == 0 {
    commands := strings.SplitN(*rawMessage, " ", 2)

    var response string
    var broadcast string
    switch commands[0] {
    default:
      response = "No command found!"
    case "/nickname":
      broadcast = nickname[*uid] + " change nickname to " + commands[1]
      nickname[*uid] = commands[1]
    }
    if len(response) > 0 {
      ResponseMessage(ws, &response)
    }
    if len(broadcast) > 0 {
      BroadcastMessage(ws, &broadcast)
    }
    return
  }

  message := (&chat.Message{Name: nickname[*uid], Content: *rawMessage}).String()
  BroadcastMessage(ws, &message)
}

func wsHandler(ws *websocket.Conn) {
  var err error
  var clientMessage string

  defer ws.Close()

  uid := fmt.Sprintf("%d-%d", clients.Len(), time.Now().Unix())
  nickname[uid] = "Guest " + uid

  item := clients.PushBack(ws)
  defer func() {
    clients.Remove(item)
    fmt.Printf("Current online clients: %d\n", clients.Len())
    delete(nickname, uid)
  }()



  for {
    if err = websocket.Message.Receive(ws, &clientMessage); err != nil {
      return
    }
    ParseCommand(&clientMessage, ws, &uid)
  }

}

func main() {
  clients = list.New()
  nickname = make(map[string]string)

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
