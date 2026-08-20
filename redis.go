package main
import(
	"context"
	"log"
	"github.com/redis/go-redis/v9"
)
var ctx=context.Background()
var rdb *redis.Client

func connectRedis(){
	rdb=redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err:=rdb.Ping(ctx).Result()
	if err!=nil{
		log.Fatal("failed to connect to Redis", err)
	}
	log.Println("Connected to Redis")
}