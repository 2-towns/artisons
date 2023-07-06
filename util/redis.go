package util

import (
	"github.com/redis/go-redis/v9"
)

// RedisClient is the client to use for Redis interactions
var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
