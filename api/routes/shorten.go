package routes

import (
	"time"
	"github.com/ubayed-bin-sufian/shorten-url-fiber-redis/database"
	"os"
)

type request struct {
	URL         string 			`json:"url"`
	CustomShort string 			`json:"short"`
	Expiry      time.Duration 	`json:"expiry"`
}

type response struct {
	URL          		string			`json:"url"`
	CustomShort  		string			`json:"short"`
	Expiry       		time.Duration	`json:"expiry"`
	XRateRemaining 		int				`json:"rate-limit"`
	XRateLimitReset  	time.Duration	`json:"rate-limit-reset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// Implement rate limiting
	// Check if the IP address of the user is already entered in the redis DB
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.IP().Result())
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// Check if the input sent by the user is an actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// Check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "You can't hack the system (:",
		})
	}

	// Enforce https, SSL
	body.URL = helpers.EnforceSSL(body.URL)

	r2.Decr(database.Ctx, c.IP())
}