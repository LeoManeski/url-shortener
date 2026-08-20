package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}

type ShortenRequest struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"`
	Tags []string `json:"tags"`
	Description string `json:"description"`
}

func shortenURL(c *fiber.Ctx) error {
	var req ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "url is required"})
	}

	code := generateCode(6)

	rdb.HSet(ctx, "link:"+code, map[string]interface{}{
			"url": req.URL,
			"clicks": 0,
			"created": time.Now().Format("2006-01-02 15:04:05"),
			"description": req.Description,
		})
	
	embedding, err:=getEmbedding(req.Description)
	
	if err!=nil{
		fmt.Println("embedding error:", err)
	}else{
		fmt.Printf("embedding saved, length: %d bytes\n", len(embedding))
		rdb.Do(ctx, "HSET", "link:"+code, "embedding", embedding)
	}

	if req.TTL > 0 {
		rdb.Expire(ctx, "link:"+code, time.Duration(req.TTL)*time.Second)
	} 

	rdb.LPush(ctx, "recent:links", code)
	rdb.LTrim(ctx, "recent:links", 0, 9)

	rdb.XAdd(ctx, &redis.XAddArgs{
		Stream:"events",
		Values: map[string]interface{}{
			"type":"link_created",
			"code": code,
			"url": req.URL,
		},
	})


	for _, tag := range req.Tags{
		rdb.SAdd(ctx, "tag:"+tag, code)
	}

	return c.Status(201).JSON(fiber.Map{
		"short_url": "http://localhost:3000/" + code,
		"code":      code,
	})
}
func getRecent(c *fiber.Ctx) error{
	links, err:=rdb.LRange(ctx, "recent:links", 0, -1).Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "couldn't fetch recent links"})
	}

	return c.JSON(fiber.Map{
		"recent":links,
	})
}
func getByTag(c *fiber.Ctx) error{
	tag:=c.Params("tag")

	links, err:=rdb.SMembers(ctx, "tag:"+tag).Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "couldn't find any tags"})
	}
	return c.JSON(fiber.Map{
		"tag": links,
	})
}
func getTop(c *fiber.Ctx) error{
	links, err:=rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 9).Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "couldn't open leaderboard"})
	}
	return c.JSON(fiber.Map{
		"top":links,
	})
}

func getEvents(c *fiber.Ctx) error{
	events, err:=rdb.XRevRange(ctx, "events", "+", "-").Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "coudln't fetch events"})
	}
	return c.JSON(fiber.Map{
		"events": events,
	})
}

func getActiveDays(c *fiber.Ctx) error{
	code:=c.Params("code")
	days, err:=rdb.BitCount(ctx, "active:"+code, nil).Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "couldnt get days"})
	}

	return c.JSON(fiber.Map{
		"code": code,
		"active_days": days,
	})
}

func getLocations(c *fiber.Ctx) error{
	code:=c.Params("code")

	locations, err:=rdb.GeoSearch(ctx, "geo:"+code, &redis.GeoSearchQuery{
		Longitude: 21.4314,
		Latitude: 41.9965,
		Radius: 10000,
		RadiusUnit: "km",
	}).Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "couldnt fetch locations"})
	}
	return c.JSON(fiber.Map{
		"code": code,
		"locations": locations,
	})
}

func searchLinks(c *fiber.Ctx) error{
	query:=c.Query("q")
	if query == ""{
		return c.Status(400).JSON(fiber.Map{"error": "query is required"})
	}

	queryEmbedding, err:=getEmbedding(query)
	if err!=nil{
		return c.Status(404).JSON(fiber.Map{"error": "embedding failed"})
	}
	
	raw, err:=rdb.Do(ctx, "FT.SEARCH", "idx:links",
		"*=>[KNN 3 @embedding $vec AS score]",
		"PARAMS", "2", "vec", queryEmbedding,
		"RETURN", "3", "url", "description", "score",
		"SORTBY", "score", "ASC",
		"DIALECT", "2",
	).Result()

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	results := []fiber.Map{}
	rawMap, ok := raw.(map[interface{}]interface{})
	if !ok {
		return c.JSON(fiber.Map{"query": query, "results": results})
	}

	rawResults, ok := rawMap["results"].([]interface{})
	if !ok {
		return c.JSON(fiber.Map{"query": query, "results": results})
	}

	for _, item := range rawResults {
		itemMap, ok := item.(map[interface{}]interface{})
		if !ok {
			continue
		}

		entry := fiber.Map{}

		if id, ok := itemMap["id"].(string); ok {
			entry["key"] = id
		}

		if attrs, ok := itemMap["extra_attributes"].(map[interface{}]interface{}); ok {
			for k, v := range attrs {
				if ks, ok := k.(string); ok {
					if vs, ok := v.(string); ok {
						entry[ks] = vs
					}
				}
			}
		}

		results = append(results, entry)
	}

	return c.JSON(fiber.Map{
		"query": query, 
		"results": results,
	})
}