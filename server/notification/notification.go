package notification

import (
	"context"
	"github.com/dapr/go-sdk/client"
)

type Notification struct {
	Client     client.Client
	PubSubName string
	TopicName  string
}

func NewNotification(pubSubName string, topicName string) (notification *Notification, err error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	return &Notification{
		Client:     daprClient,
		PubSubName: pubSubName,
		TopicName:  topicName,
	}, nil
}

func (n *Notification) SendNotification(message []byte) error {
	return n.Client.PublishEvent(context.Background(), n.PubSubName, n.TopicName, message)
}
