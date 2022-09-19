package rds

import (
	"animeGAN2go/plugins"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     plugins.RedisHost,
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Send(name string, structure map[string]string) {

	uid := uuid.New().String()

	fmt.Println("send task with uuid - " + uid)

	rdb.HMSet(ctx, uid, structure)

	rdb.LPush(ctx, name, uid)
}

func Receive(name string) map[string]string {
	uid, err := rdb.RPop(ctx, name).Result()

	var data map[string]string

	if err == redis.Nil {
		fmt.Println("no jobs ")
		return data
	}

	data, err = rdb.HGetAll(ctx, uid).Result()

	if err != nil {
		fmt.Println(err)
		fmt.Println("something wrong")
	}

	return data

}
