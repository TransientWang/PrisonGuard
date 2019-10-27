package src

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
)



func TestName(t *testing.T) {
	exit := make(chan os.Signal)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	scoreRedis := RedisKeyManager{
		Rds:        client,
		ManagerKey: "key-manger",
		KeyGenerateStrategy: func(pattern string) (result []*redis.Z) {
			result = append(result, &redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: pattern,
			})
			return
		},
		SurvivalMinute: 1,
		Error:          make(chan error),
	}
	fmt.Println(scoreRedis)
	go scoreRedis.BeginWork()
	go func() {
		for _, key := range []string{"a", "b", "c"} {
			fmt.Println(client.HSet(key, key, key))
			scoreRedis.AddKey(key)
		}
	}()
	exit <- os.Interrupt
}

func TsestNames(t *testing.T) {
	var err chan int
	err = make(chan int)
	go func() {
		fmt.Println("begin")
		err <- 1
		fmt.Println("end")
	}()
	go func() {
		fmt.Println("beginc")
		a := <-err
		fmt.Println(a)
	}()
	time.Sleep(3 * time.Second)
}
