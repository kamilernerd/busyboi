package main

import (
	"context"
	"encoding/json"
	"time"

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

func worker(rawConfig []byte, ctx context.Context, bb *Busyboi) {
	var conf JobConfig
	err := json.Unmarshal(rawConfig, &conf)

	if err != nil {
		// Raport error
		return
	}

	cache := bb.cache.Has(conf.Url)

	// Cache hit!
	if cache != nil {
		// fmt.Println(bb.cache.table)
		return
	}

	// Re-queue job is site is unreachable
	ok := canReach(conf.Url)
	if !ok {
		bb.mq.RabbitMqAddMessages(conf)
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

	data := parse(conf.Fields, DOM, "")

	bb.cache.AddOrUpdate(Bucket{
		Url:       conf.Url,
		Refreshed: time.Now(),
		Data:      data,
	})
}
