// Package VAPID allows webpush notifications.
// It's called vapid to avoid confusion with webpush
// go package.
package vapid

import (
	"context"
	"encoding/json"
	"gifthub/conf"
	"log/slog"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// Send a push notification on a device.
func Send(ctx context.Context, device, message string) error {
	l := slog.With(slog.String("message", message))
	l.LogAttrs(ctx, slog.LevelInfo, "sending a new notification")

	s := &webpush.Subscription{}
	json.Unmarshal([]byte(device), s)

	l.LogAttrs(ctx, slog.LevelInfo, "preparing to send a push notification")

	resp, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
		Subscriber:      conf.VapidEmail,
		VAPIDPublicKey:  conf.VapidPublicKey,
		VAPIDPrivateKey: conf.VapidPrivateKey,
		TTL:             30,
	})

	if err != nil {
		l.LogAttrs(ctx, slog.LevelError, "cannot send the notification", slog.String("err", err.Error()), slog.Int("status", resp.StatusCode))
	} else {
		l.LogAttrs(ctx, slog.LevelInfo, "notification sent", slog.Int("status", resp.StatusCode))
	}

	defer resp.Body.Close()

	return err
}
