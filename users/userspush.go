package users

import (
	"context"
	"errors"
	"fmt"
	"gifthub/db"
	"log"
)

// AddWPToken registers a vapid webpush token
// to receive push notifications.
// The token is the string representation of the
// JSON token.
func (u User) AddWPToken(token string) error {
	if token == "" {
		log.Printf("input_validation_fail the token is required")
		return errors.New("user_wptoken_required")
	}

	ctx := context.Background()
	if _, err := db.Redis.HSet(ctx, fmt.Sprintf("user:%d", u.ID), "wptoken:"+u.SID, token).Result(); err != nil {
		log.Printf("ERROR: sequence_fail when storing the wptoken %s : %s", token, err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}

// DeleteWPToken removes a vapid webpush token linked to
// a session.
func (u User) DeleteWPToken(sid string) error {
	if sid == "" {
		log.Printf("input_validation_fail the session id is required")
		return errors.New("unauthorized")
	}

	ctx := context.Background()
	if _, err := db.Redis.HDel(ctx, fmt.Sprintf("user:%d", u.ID), "wptoken:"+u.SID).Result(); err != nil {
		log.Printf("ERROR: sequence_fail when deleting the wptoken for the session %s : %s", u.SID, err.Error())
		return errors.New("something_went_wrong")
	}

	return nil
}
