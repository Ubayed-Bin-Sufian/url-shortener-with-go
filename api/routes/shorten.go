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