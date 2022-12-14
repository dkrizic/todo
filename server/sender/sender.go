package sender

import (
	"context"
	"github.com/dapr/go-sdk/client"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	client     client.Client
	PubSubName string
	TopicName  string
}

func NewSender(pubSubName string, topicName string) (notification *Sender, err error) {
	daprClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	return &Sender{
		client:     daprClient,
		PubSubName: pubSubName,
		TopicName:  topicName,
	}, nil
}

func (n *Sender) SendNotification(message []byte) error {
	llog := log.WithFields(log.Fields{
		"pubsubName": n.PubSubName,
		"topicName":  n.TopicName,
		"message":    string(message),
	})
	llog.Debug("Sending sender")
	err := n.client.PublishEvent(context.Background(), n.PubSubName, n.TopicName, message)
	if err != nil {
		llog.WithError(err).Warn("Unable to send sender")
		return err
	}
	return nil
}
