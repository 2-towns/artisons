// Package VAPID allows webpush notifications.
// It's called vapid to avoid confusion with webpush
// go package.
package vapid

import (
	"encoding/json"
	"gifthub/conf"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// Send a push notification on a device.
func Send(device, message string) error {
	s := &webpush.Subscription{}
	json.Unmarshal([]byte(device), s)

	resp, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
		Subscriber:      conf.VapidEmail,
		VAPIDPublicKey:  conf.VapidPublicKey,
		VAPIDPrivateKey: conf.VapidPrivateKey,
		TTL:             30,
	})

	defer resp.Body.Close()

	return err
}
