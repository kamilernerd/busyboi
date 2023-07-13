package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestWorker(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	bb := &Crawlerbox{
		queueMsgs:      make(chan amqp.Delivery),
		wg:             sync.WaitGroup{},
		workersCounter: 1,
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

	for m := range bb.queueMsgs {
		t := time.Now()
		timeoutWorker := time.Duration(math.Ceil(1e9 / (ExecutionsPerWorker / float64(bb.workersCounter))))

		timeoutUntilNext := -(time.Since(t) - timeoutWorker)

		// Timeouts bigger than 1sec are no go...
		if timeoutUntilNext > time.Second * 1 {
			timeoutUntilNext = time.Microsecond * 100
		}

		if timeoutUntilNext > 0 {
			time.Sleep(timeoutUntilNext)
		}
		go worker(m.Body, ctx, bb)
	}

	bb.wg.Wait()
	close(bb.queueMsgs)
}

func TestThrottle(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	urls := []string{
		"http://localhost:8080",
	}

	for i := 0; i < 50; i++ {
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

	}
}

func TestTestServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testserver()
}
