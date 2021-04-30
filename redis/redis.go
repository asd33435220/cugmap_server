package redis

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)
// key规则： 地点 p+placeCode
var Rdb *redis.Client

func init() {
	initClient()
}

func initClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func Save(){
	
}
