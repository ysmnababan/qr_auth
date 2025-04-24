package config

import (
	"log"
	"qr_auth/pusherutil"
	"qr_auth/redisutil"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
)

type Config struct {
	PusherClient *pusher.Client
	Redis        *redis.Client
}

var Cfg Config

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	Cfg = Config{
		PusherClient: pusherutil.NewPusherClient(),
		Redis:        redisutil.NewRedisClient("localhost:6379", ""),
	}
}
