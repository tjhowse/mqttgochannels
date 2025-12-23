package mqttgochannels

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTgoChannels struct {
	mqtt.Client
	topicChannels map[string]chan mqtt.Message
}

func NewMQTTgoChannels(client mqtt.Client) *MQTTgoChannels {
	return &MQTTgoChannels{
		Client:        client,
		topicChannels: make(map[string]chan mqtt.Message),
	}
}

func (m *MQTTgoChannels) SubscribeGetChannel(topic string, qos byte) (chan mqtt.Message, error) {
	// Subscribes to a topic and returns a channel to receive messages

	// Check if we already have a channel for this topic, if so: return it.
	if ch, ok := m.topicChannels[topic]; ok {
		return ch, nil
	}

	// If not create a new channel and subscribe to the topic
	ch := make(chan mqtt.Message)
	token := m.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		// Inline function to dispatch the subscription responses to the correct channel
		if ch, ok := m.topicChannels[topic]; ok {
			ch <- msg
			return
		}
	})

	// Block until the subscription is complete
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	// Store the channel in the map and return it
	m.topicChannels[topic] = ch
	return ch, nil
}
