package redilock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	//透過lua腳本執行，達到原子性操作
	luaDelete = `
			if redis.call('GET', KEYS[1]) == ARGV[1] then
				return redis.call('DEL', KEYS[1])
			else
				return 0
			end
		`
)

var (
	Rediclient *redis.Client
)

//InitRedlock : redis cleint 設定
func InitRedlock(rediclient *redis.Client) {
	Rediclient = rediclient
}

type RedLock struct {
	key        string
	expiration time.Duration

	value string
}

func NewRediLock(key string, exp time.Duration) *RedLock {
	return &RedLock{
		key:        key,
		expiration: exp,
	}
}

func (rlock *RedLock) RediLock(ctx context.Context) (bool, error) {

	// 隨機值 uuid. or used snowflake.
	rlock.value = uuid.New().String()

	//設定key值並做retry機制
	{
		retry := 0

		for {
			set, err := Rediclient.SetNX(ctx, rlock.key, rlock.value, rlock.expiration).Result()

			if err != nil {
				panic(err.Error())
			}
			if set == true {
				return true, nil
			}

			if retry >= 5 {
				return false, nil
			}
			retry++
			time.Sleep(time.Millisecond * 50)
		}
	}

}

func (rlock *RedLock) RediUnLock(ctx context.Context) (bool, error) {

	result, err := Rediclient.Eval(ctx, luaDelete, []string{rlock.key}, []string{rlock.value}).Result()
	if err != nil {
		return false, err
	}
	if result != int64(1) {
		return false, err
	}
	return true, nil
}
