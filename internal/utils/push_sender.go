package utils

import (
	"chechnya-product/internal/models"
	"encoding/json"
	"github.com/SherClockHolmes/webpush-go"
	"os"
)

var (
	VAPIDPublicKey  = os.Getenv("VAPID_PUBLIC_KEY")
	VAPIDPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
)

func SendPush(sub models.SubscribeRequest, title, body string) error {
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.Keys.P256dh,
			Auth:   sub.Keys.Auth,
		},
	}

	payload, _ := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})

	resp, err := webpush.SendNotification(payload, subscription, &webpush.Options{
		Subscriber:      "mailto:test@example.com",
		VAPIDPublicKey:  VAPIDPublicKey,
		VAPIDPrivateKey: VAPIDPrivateKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
