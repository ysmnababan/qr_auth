package pusherutil

import (
	"os"

	"github.com/pusher/pusher-http-go/v5"
)

var QR_CHANNEL = "qr_channel"
var QR_EVENT = "qr_event"

type IPusher interface {
	Trigger(channel string, eventName string, data any) error
}

type PusherClient struct {
	client *pusher.Client
}

func NewPusherClient() *PusherClient {
	return &PusherClient{
		client: &pusher.Client{
			AppID:   os.Getenv("PUSHER_APP_ID"),
			Key:     os.Getenv("PUSHER_KEY"),
			Secret:  os.Getenv("PUSHER_SECRET"),
			Cluster: os.Getenv("PUSHER_CLUSTER"),
			Secure:  true,
		},
	}
}

func (p *PusherClient) Trigger(channel string, eventName string, data any) error {
	return p.client.Trigger(channel, eventName, data)
}
