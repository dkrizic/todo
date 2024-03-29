package sender

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type Sender struct {
	PubSubName string
	TopicName  string
}

func NewSender(pubSubName string, topicName string) (notification *Sender, err error) {
	return &Sender{
		PubSubName: pubSubName,
		TopicName:  topicName,
	}, nil
}

func (n *Sender) SendNotification(ctx context.Context, message []byte) error {
	ctx, span := otel.Tracer("sender").Start(ctx, "SendNotification")
	defer span.End()

	llog := log.WithFields(log.Fields{
		"pubsubName": n.PubSubName,
		"topicName":  n.TopicName,
		"message":    string(message),
	})
	llog.Debug("Sending sender")

	client, err := dapr.NewClient()
	if err != nil {
		llog.WithError(err).Warn("Unable to create dapr client")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	defer client.Close()

	err = client.PublishEvent(ctx, n.PubSubName, n.TopicName, message)
	if err != nil {
		llog.WithError(err).Warn("Unable to send sender")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	return nil
}
