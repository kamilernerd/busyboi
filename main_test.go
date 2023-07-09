package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestWorker(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	bb := &Busyboi{
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
						},
						{
							Name:     "Image 2",
							Selector: "body .product_image",
							Attr:     "src",
						},
					},
				},
				{
					Name:     "Title",
					Selector: "body .product_title",
				},
				{
					Name:     "Price",
					Selector: "body .product_price",
					Filters: []string{
						"number",
					},
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
		chromedp.UserAgent("Busyboi"),
		chromedp.WindowSize(1920, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	i := 0
	for m := range bb.queueMsgs {
		go worker(m.Body, ctx)
		i++
	}
}
