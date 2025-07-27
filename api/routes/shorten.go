package routes

import (
	"time"
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

}