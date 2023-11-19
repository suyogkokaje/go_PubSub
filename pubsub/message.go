package pubsub

type Message struct {
	topic string
	body string
}

func NewMessage(msg string, topic string) (* Message) {
	return &Message{
		topic: topic,
		body: msg,
	}
}
func (m *Message) GetTopic() string {
	return m.topic
}
func (m *Message) GetMessageBody() string {
	return m.body
}