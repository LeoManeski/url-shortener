package main

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)
func redirectURL(c *fiber.Ctx) error{
	code:=c.Params("code")
	

	// if err!=nil{
	// 	return c.Status(404).JSON(fiber.Map{"error": "link not found"})
	// }

	// pipe:=rdb.TxPipeline()
	// pipe.HIncrBy(ctx, "link:"+code, "clicks", 1)
	// pipe.ZIncrBy(ctx, "leaderboard", 1, code)
	// _, err=pipe.Exec(ctx)

	script:=redis.NewScript(`
		local exists=redis.call('EXISTS', KEYS[1])
		if exists==0 then
			return nil
		end
		local url = redis.call('HGET', KEYS[1], 'url')
		redis.call('HINCRBY', KEYS[1], 'clicks', 1)
		redis.call('ZINCRBY', KEYS[2], 1, ARGV[1])
		return url
	`)
	result, err:=script.Run(ctx, rdb, []string{"link:"+code, "leaderboard"}, code).Result()
	if err!=nil{
		return c.Status(404).JSON(fiber.Map{"error": "link not found"})
	}
	url := result.(string)

	rdb.PFAdd(ctx, "visitors:"+code, c.IP())

	dayOfYear:=time.Now().YearDay()
	rdb.SetBit(ctx, "active:"+code, int64(dayOfYear), 1)

	rdb.GeoAdd(ctx, "geo:"+code, &redis.GeoLocation{
		Name:c.IP(),
		Longitude: 21.4314,
		Latitude: 41.9965,
	})

	rdb.Publish(ctx, "link:clicks", code)

	rdb.XAdd(ctx, &redis.XAddArgs{	
		Stream:"events",
		Values:map[string]interface{}{
			"type":"link_clicked",
			"code": code,
		},
	})

	return c.Redirect(url)
}