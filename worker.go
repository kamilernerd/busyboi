package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
)

type JobConfig struct {
	Collection string           `json:"name"`
	Url        string           `json:"url"`
	Fields     []JobConfigField `json:"fields"`
	Id         string           `json:"id,omitempty"`
}

type JobConfigField struct {
	Name     string           `json:"name"`
	Selector string           `json:"selector"`
	Attr     string           `json:"attr,omitempty"`
	Children []JobConfigField `json:"children,omitempty"`
	Regex    string           `json:"regex,omitempty"`
}

func worker(rawConfig []byte, ctx context.Context) {
	var conf JobConfig
	err := json.Unmarshal(rawConfig, &conf)

	if err != nil {
		// Raport error
		return
	}

	// Re-queue job is site is unreachable
	ok := canReach(conf.Url)
	if !ok {
		RabbitMqAddMessages(conf)
		return
	}

	// Check if we can crawl
	ok = robots(conf.Url)
	if !ok {
		return
	}

	c, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var DOM string

	// run task list
	err = chromedp.Run(c,
		chromedp.Navigate(conf.Url),
		chromedp.InnerHTML("html", &DOM),
	)

	if err != nil {
		return
	}

	fmt.Print(parse(conf.Fields, DOM, ""))
}
