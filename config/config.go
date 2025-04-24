package config

import (
	"log"
	"qr_auth/pusherutil"
	"qr_auth/redisutil"

	"github.com/joho/godotenv"
)

type Config struct {
	PusherClient *pusherutil.PusherClient
	Redis        *redisutil.RedisClient
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
