package main

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestWorker(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	bb := &Crawlerbox{
		queueMsgs: make(chan amqp.Delivery, 100),
	}

	go testserver()

	urls := []string{
		"http://localhost:8080",
		"http://localhost:1337",
	}

	for i, k := range urls {
		RabbitMqAddMessages(JobConfig{
			Collection: "products",
			Url:        k,
			Fields: []JobConfigField{
				{
					Name:     "Images",
					Selector: ".images",
					Attr:     "src",
					Children: []JobConfigField{
						{
							Name:     "Image 1",
							Selector: "body .product_image",
							Attr:     "src",
							Regex:    "(\\.jpg)",
						},
						{
							Name:     "Image 2",
							Selector: "body .product_image",
							Attr:     "src",
						},
					},
					Regex: "([a-z]+\\.jpg)",
				},
				{
					Name:     "Title",
					Selector: "body .product_title",
				},
				{
					Name:     "Price",
					Selector: "body .product_price",
					Regex:    "(\\.[0-9]+)",
				},
				{
					Name:     "Some url",
					Selector: "body a[href]",
					Attr:     "href",
				},
			},
		})

		fmt.Printf("Added job # %d\n", i)
	}

	go RabbitMqGetMessages(bb)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent("Crawlerbox"),
		chromedp.WindowSize(1920, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	var wg sync.WaitGroup

	for m := range bb.queueMsgs {
		go worker(m.Body, ctx, &wg)
	}

	wg.Wait()
	close(bb.queueMsgs)
}
