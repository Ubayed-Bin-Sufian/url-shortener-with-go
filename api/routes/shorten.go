package routes

import (
	"os"
	"time"
	"strconv"
	"github.com/ubayed-bin-sufian/url-shortener-with-go/api/database"
	"github.com/ubayed-bin-sufian/url-shortener-with-go/api/helpers"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
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
		val, _ = r2.Get(database.Ctx, c.IP()).Result()
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

	// Check if the user has sent us a custom short URL
	// Check if the custom short URL is already taken
	// If not, create a new short URL with random ID using UUID
	
	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom short URL is already taken",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to the server",
		})
	}

	// Create response from the function
	resp := response {
		URL:				body.URL,
		CustomShort:		"",
		Expiry:				body.Expiry,
		XRateRemaining:		10,
		XRateLimitReset: 	30, 
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	// ttl value from DB
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}