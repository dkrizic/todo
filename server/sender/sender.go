package sender

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	log "github.com/sirupsen/logrus"
	"time"
)

type Sender struct {
	client     dapr.Client
	PubSubName string
	TopicName  string
}

func NewSender(pubSubName string, topicName string) (notification *Sender, err error) {
	var daprClient dapr.Client
	// try 10 times to connect to dapr and wait 1 seconds after each try
	max := 10
	for i := 0; i < max; i++ {
		daprClient, err = dapr.NewClient()
		if err == nil {
			break
		}
		log.WithFields(log.Fields{
			"try": i,
			"max": max,
		}).WithError(err).Warn("Unable to connect to dapr")
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return &Sender{
		client:     daprClient,
		PubSubName: pubSubName,
		TopicName:  topicName,
	}, nil
}

func (n *Sender) SendNotification(ctx context.Context, message []byte) error {
	llog := log.WithFields(log.Fields{
		"pubsubName": n.PubSubName,
		"topicName":  n.TopicName,
		"message":    string(message),
	})
	llog.Debug("Sending sender")
	err := n.client.PublishEvent(ctx, n.PubSubName, n.TopicName, message)
	if err != nil {
		llog.WithError(err).Warn("Unable to send sender")
		return err
	}
	return nil
}
