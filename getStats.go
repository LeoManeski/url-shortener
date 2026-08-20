package main

import (
	"github.com/gofiber/fiber/v2"
)
func getStats(c *fiber.Ctx)error{
	code:=c.Params("code")

	pipe:=rdb.Pipeline()
	dataCmd:=pipe.HGetAll(ctx, "link:"+code)
	ttlCmd:=pipe.TTL(ctx, "link:"+code)
	uniqueCmd:=pipe.PFCount(ctx, "visitors:"+code)
	_, err:=pipe.Exec(ctx)

	if err!=nil{
		return c.Status(404).JSON(fiber.Map{"error": "Link not found"})
	}


	return c.JSON(fiber.Map{
		"url": dataCmd.Val()["url"],
		"clicks": dataCmd.Val()["clicks"],
		"created": dataCmd.Val()["created"],
		"ttl": ttlCmd.Val(),
		"unique_visitors": uniqueCmd.Val(),
	})
}