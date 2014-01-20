package chatter

type Message struct {
  Name string
  Content string
}

func (m *Message) String() string {
  return m.Name + ":" + m.Content
}

func NewMessage(name string, content string) *Message {
  return &Message{Name: name, Content: content}
}
