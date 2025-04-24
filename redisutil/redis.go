package redisutil

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var REDIS_QR_LOGIN_PREFIX = "qrlogin:"

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

type ICache interface {
	Set(key string, val any, exp time.Duration) (err error)
	Get(key string) (val string, err error)
	Delete(key string) (delCount int64, err error)
}

func NewRedisClient(host string, password string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
	redisClient := new(RedisClient)
	redisClient.client = client
	redisClient.ctx = ctx
	return redisClient
}

func (r *RedisClient) Set(key string, val any, exp time.Duration) (err error) {
	return r.client.Set(r.ctx, key, val, exp).Err()
}

func (r *RedisClient) Get(key string) (val string, err error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Delete(key string) (delCount int64, err error) {
	return r.client.Del(r.ctx, key).Result()
}
