package main
import(
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)
func main(){
	connectRedis()
	createIndex()
	app:=fiber.New()	
	app.Use(cors.New())

	app.Post("/shorten", shortenURL)
	app.Get("/stats/:code", getStats)
	app.Get("/recent", getRecent)
	app.Get("/tags/:tag", getByTag)
	app.Get("/top", getTop)
	app.Get("/events", getEvents)
	app.Get("/search", searchLinks)

	
	app.Get("/:code", redirectURL)
	app.Get("/active/:code", getActiveDays)
	app.Get("/locations/:code", getLocations)
	go startSubscriber()

	log.Fatal(app.Listen(":3000"))
}