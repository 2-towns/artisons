// Package db provides redis storage
package db

import (
	"github.com/redis/go-redis/v9"
)

// Client is the client to use for Redis interactions
var Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
