// When a user uses the short url, we check the redis database for the short url 
// If it exists, we redirect to the original URL

package routes

import (
	"github.com/ubayed-bin-sufian/url-shortener-with-go/api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
)

func ResolveURL(c *fiber.Ctx) error {

	url := c.Params("url")

	r := database.CreateClient()
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short URL not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error retrieving short URL from the database",
		})
	}
	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}