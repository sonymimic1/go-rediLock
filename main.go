package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	redilock "github.com/sonymimic1/go-rediLock/redilock"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var ctx, cancel = context.WithCancel(context.Background())

	defer func() {
		cancel()
	}()

	redilock.InitRedlock(rdb)

	redlock := redilock.NewRediLock("LOCK:KEY", time.Second*20)

	if ok, err := redlock.RediLock(ctx); !ok {
		if err != nil {
			fmt.Printf("操作失敗 \n")
			return
		}
		fmt.Printf("Lock稍後再嘗試，有其他SERVICE佔用著Redis KEY \n")
		return
	}
	fmt.Printf("LOCK OK! \n")

	// DoSometing...
	time.Sleep(time.Second * 15)

	if ok, err := redlock.RediUnLock(ctx); !ok {
		if err != nil {
			fmt.Printf("UnLock Redis操作失敗 \n")
			return
		}
		fmt.Printf("UnLock稍後再嘗試，有其他SERVICE佔用著Redis KEY \n")
		return
	}

	fmt.Printf("UnLOCK OK! \n")
}
