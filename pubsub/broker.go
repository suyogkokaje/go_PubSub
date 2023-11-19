package pubsub

import (
	"fmt"
	"sync"
)

type Subscribers map[string]*Subscriber

type Broker struct {
	subscribers Subscribers 
	topics map[string]Subscribers 
	mut sync.RWMutex 
}

func NewBroker() (*Broker){
	return &Broker{
		subscribers: Subscribers{},
		topics: map[string]Subscribers{},
	}
}

func (b *Broker)AddSubscriber()(*Subscriber){
	b.mut.Lock()
	defer b.mut.Unlock()
	id, s := CreateNewSubscriber()
	b.subscribers[id] = s;
	return s
}

func (b *Broker)RemoveSubscriber(s *Subscriber)(){
	for topic := range(s.topics){
		b.Unsubscribe(s, topic)
	}
	b.mut.Lock()
	delete(b.subscribers, s.id)
	b.mut.Unlock()
	s.Destruct()
}

func (b *Broker)Broadcast(msg string, topics []string){
	for _, topic:=range(topics) {
		for _, s := range(b.topics[topic]){
			m:= NewMessage(msg, topic)
			go (func(s *Subscriber){
				s.Signal(m)
			})(s)
		}
	}
}

func (b *Broker) GetSubscribers(topic string) int {
	b.mut.RLock()
	defer b.mut.RUnlock()
	return len(b.topics[topic])
}

func (b *Broker) Subscribe(s *Subscriber, topic string) {
	b.mut.Lock()
	defer b.mut.Unlock()

	if  b.topics[topic] == nil {
		b.topics[topic] = Subscribers{}
	}
	s.AddTopic(topic)
	b.topics[topic][s.id] = s
	fmt.Printf("%s Subscribed for topic: %s\n", s.id, topic)
}

func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	b.mut.RLock()
	defer b.mut.RUnlock()

	delete(b.topics[topic], s.id)
	s.RemoveTopic(topic)
	fmt.Printf("%s Unsubscribed for topic: %s\n", s.id, topic)
}

func (b *Broker) Publish(topic string, msg string) {
	b.mut.RLock()
	bTopics := b.topics[topic]
	b.mut.RUnlock()
	for _, s := range bTopics {
		m:= NewMessage(msg, topic)
		if !s.active{
			return
		}
		go (func(s *Subscriber){
			s.Signal(m)
		})(s)
	}
}
