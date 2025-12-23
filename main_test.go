package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fullPath := fmt.Sprintf("tcp://%s:%s", os.Getenv("MQTT_BROKER_HOSTNAME"), os.Getenv("MQTT_BROKER_PORT"))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fullPath)
	opts.SetClientID("testing")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	// opts.SetDefaultPublishHandler(f)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v. Attempting to reconnect...", err)
		for {
			// Attempt to reconnect every 5 seconds. Not sure if this will work well.
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Printf("Reconnection failed: %v. Retrying in 5 seconds...", token.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			log.Println("Reconnected to MQTT broker.")
			break
		}

	})

	client := NewMQTTgoChannels(mqtt.NewClient(opts))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Error connecting to MQTT broker:", token.Error())
	}

	msgChannel, err := client.SubscribeGetChannel("msh/ANZ/2/json/MediumFast/#", 0)
	if err != nil {
		fmt.Println("Error subscribing to topic:", err)
		return
	}

	go func() {
		for msg := range msgChannel {
			fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
		}
	}()

	select {}

}

// Received message on topic msh/ANZ/2/json/MediumFast/!849a88e4:
// {
//     "channel": 0,
//     "from": 2224720100,
//     "hop_start": 5,
//     "hops_away": 0,
//     "id": 2474236603,
//     "payload": {
//         "text": "Over too fast here. I hope that's not all we get."
//     },
//     "sender": "!849a88e4",
//     "timestamp": 1766452827,
//     "to": 4294967295,
//     "type": "text"
// }
