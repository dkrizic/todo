package notification

import "github.com/dapr/go-sdk/client"

type Notification struct {
	Client client.Client
}

func NewNotification(client client.Client) *Notification {
	return &Notification{
		Client: client,
	}
}

func (n *Notification) SendNotification() error {
	return n.Client.PublishEvent("pubsub", "notification", []byte("Hello World!"))
}
