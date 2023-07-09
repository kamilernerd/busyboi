package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

/*
 * Send a GET request to check if url is reachable
 */
func canReach(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}

	return resp.StatusCode == 200
}

/*
 * Sent a GET request and try to locate robots.txt
 * When found check if our UserAgent is allowed to access requested path
 * also check if we can access /
 */
func robots(rawUrl string) bool {
	u, _ := url.Parse(rawUrl)
	resp, err := http.Get(fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host))
	if err != nil {
		return false
	}

	// Check if allowed
	robots, err := robotstxt.FromResponse(resp)
	group := robots.FindGroup("Busyboi")

	can := true
	if !group.Test(u.Path) {
		can = false
	}

	if !group.Test("/") {
		can = false
	}

	return can
}
