package pubsub

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
)


type Subscriber struct {
	id string 
	messages chan* Message
	topics map[string]bool 
	active bool 
	mutex sync.RWMutex
}

func CreateNewSubscriber() (string, *Subscriber) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%X-%X", b[0:4], b[4:8])
	return id, &Subscriber{
		id: id,
		messages: make(chan *Message),
		topics: map[string]bool{},
		active: true,
	}
}

func (s * Subscriber)AddTopic(topic string)(){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.topics[topic] = true
}

func (s * Subscriber)RemoveTopic(topic string)(){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.topics, topic)
}

func (s * Subscriber)GetTopics()([]string){
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	topics := []string{}
	for topic, _ := range s.topics {
		topics = append(topics, topic)
	}
	return topics
}

func (s *Subscriber)Destruct() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.active = false
	close(s.messages)
}

func (s *Subscriber)Signal(msg *Message) () {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.active{
		s.messages <- msg
	}
}

func (s *Subscriber)Listen() {
	for {
		if msg, ok := <- s.messages; ok {
			fmt.Printf("Subscriber %s, received: %s from topic: %s\n", s.id, msg.GetMessageBody(), msg.GetTopic())
		}
	}
}