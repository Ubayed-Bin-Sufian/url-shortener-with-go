package helpers

import (
	"os"
	"strings"
)

func EnforceSSL () string {
	if url[:4] != "http" {
		return "https://" + url
  	}
  	return url
}

func RemoveDomainError (url string) bool {

	if url == os.Getenv("DOMAIN"){
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
