package notification

import (
	"context"
	"github.com/dapr/go-sdk/client"
	log "github.com/sirupsen/logrus"
)

type Notification struct {
	Enabled    bool
	client     client.Client
	PubSubName string
	TopicName  string
}

func NewNotification(enabled bool, pubSubName string, topicName string) (notification *Notification, err error) {
	if err != nil {
		return nil, err
	}
	var daprClient client.Client
	if enabled {
		daprClient, err = client.NewClient()
		if err != nil {
			return nil, err
		}
	}
	return &Notification{
		Enabled:    enabled,
		client:     daprClient,
		PubSubName: pubSubName,
		TopicName:  topicName,
	}, nil
}

func (n *Notification) SendNotification(message []byte) error {
	if !n.Enabled {
		log.Debug("Notification is disabled")
		return nil
	}
	llog := log.WithFields(log.Fields{
		"pubsubName": n.PubSubName,
		"topicName":  n.TopicName,
		"message":    string(message),
	})
	llog.Debug("Sending notification")
	err := n.client.PublishEvent(context.Background(), n.PubSubName, n.TopicName, message)
	if err != nil {
		llog.WithError(err).Warn("Unable to send notification")
		return err
	}
	return nil
}
