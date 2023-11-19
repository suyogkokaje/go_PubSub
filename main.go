package main

import (
	"fmt"
	"math/rand"
	"time"

	"./pubsub"
)

// available topics
var availableTopics = map[string]string{
	"BTC": "BITCOIN",
	"ETH": "ETHEREUM",
	"DOT": "POLKADOT",
	"SOL": "SOLANA",
}

func pricePublisher(broker *pubsub.Broker)(){
	topicKeys := make([]string, 0, len(availableTopics))
	topicValues := make([]string, 0, len(availableTopics))

	for k, v := range availableTopics {
		topicKeys = append(topicKeys, k)
		topicValues = append(topicValues, v)
	}
	for {
		randValue := topicValues[rand.Intn(len(topicValues))] 
		msg:= fmt.Sprintf("%f", rand.Float64())
		// fmt.Printf("Publishing %s to %s topic\n", msg, randKey)
		go broker.Publish(randValue, msg)
        // Uncomment if you want to broadcast to all topics.
		// go broker.Broadcast(msg, topicValues)
		r := rand.Intn(4)
		time.Sleep(time.Duration(r) * time.Second) 
	}
}



func main(){
	broker := pubsub.NewBroker()

	s1 := broker.AddSubscriber()
    
	broker.Subscribe(s1, availableTopics["BTC"])
	broker.Subscribe(s1, availableTopics["ETH"])

    
	s2 := broker.AddSubscriber()
   
	broker.Subscribe(s2, availableTopics["ETH"])
	broker.Subscribe(s2, availableTopics["SOL"])

	go (func(){
		time.Sleep(3*time.Second)
		broker.Subscribe(s2, availableTopics["DOT"])
	})()

	go (func(){
		time.Sleep(5*time.Second)
		broker.Unsubscribe(s2, availableTopics["SOL"])
		fmt.Printf("Total subscribers for topic ETH is %v\n", broker.GetSubscribers(availableTopics["ETH"]))
	})()


	go (func(){
		time.Sleep(10*time.Second)
		broker.RemoveSubscriber(s2)
		fmt.Printf("Total subscribers for topic ETH is %v\n", broker.GetSubscribers(availableTopics["ETH"]))
	})()

	go pricePublisher(broker)
	go s1.Listen()
	go s2.Listen()

	fmt.Scanln()
	fmt.Println("Done!")
}