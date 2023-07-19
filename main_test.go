package main

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestWorker(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var queue_host string
	var queue_user string
	var queue_password string
	var queue_port string
	var queue_name string
	var concurreny_limit int
	var cache_ttl int
	flag.StringVar(&queue_host, "queue_host", "localhost", "Hostname for rabbitmq")
	flag.StringVar(&queue_user, "queue_user", "guest", "User for rabbitmq")
	flag.StringVar(&queue_password, "queue_password", "guest", "Password for rabbitmq")
	flag.StringVar(&queue_port, "queue_port", "5672", "Port for rabbitmq")
	flag.StringVar(&queue_name, "queue_name", "busyboi", "Queue name for rabbitmq")
	flag.IntVar(&concurreny_limit, "concurreny_limit", 10, "Limits how many crawling jobs can run simultaneously")
	flag.IntVar(&cache_ttl, "cache_ttl", 900, "How long should cache be stored (in seconds). Default is 15 minutes.")
	flag.Parse()

	bb := &Busyboi{
		queueMsgs: make(chan amqp.Delivery),
		mq: &Rabbitmq{
			hostname: queue_host,
			user:     queue_user,
			password: queue_password,
			port:     queue_port,
			queue:    queue_name,
		},
		cache: NewCache(time.Second * time.Duration(cache_ttl)),
	}

	go bb.cache.Invalidator()

	go testserver()

	urls := []string{
		"http://localhost:8080",
	}

	for i, k := range urls {
		bb.mq.RabbitMqAddMessages(JobConfig{
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

	go bb.mq.RabbitMqGetMessages(bb)

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

	limiter := make(chan uint, concurreny_limit)
	for m := range bb.queueMsgs {
		limiter <- 1

		go func(m amqp.Delivery) {
			worker(m.Body, ctx, bb)
			<-limiter
		}(m)
	}

	close(bb.queueMsgs)
}

func TestThrottle(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var queue_host string
	var queue_user string
	var queue_password string
	var queue_port string
	var queue_name string
	var concurreny_limit int
	var cache_ttl int
	flag.StringVar(&queue_host, "queue_host", "localhost", "Hostname for rabbitmq")
	flag.StringVar(&queue_user, "queue_user", "guest", "User for rabbitmq")
	flag.StringVar(&queue_password, "queue_password", "guest", "Password for rabbitmq")
	flag.StringVar(&queue_port, "queue_port", "5672", "Port for rabbitmq")
	flag.StringVar(&queue_name, "queue_name", "busyboi", "Queue name for rabbitmq")
	flag.IntVar(&concurreny_limit, "concurreny_limit", 10, "Limits how many crawling jobs can run simultaneously")
	flag.IntVar(&cache_ttl, "cache_ttl", 900, "How long should cache be stored (in seconds). Default is 15 minutes.")
	flag.Parse()
	bb := &Busyboi{
		queueMsgs: make(chan amqp.Delivery),
		mq: &Rabbitmq{
			hostname: queue_host,
			user:     queue_user,
			password: queue_password,
			port:     queue_port,
			queue:    queue_name,
		},
		cache: NewCache(time.Second * time.Duration(cache_ttl)),
	}

	go bb.cache.Invalidator()

	urls := []string{
		"http://localhost:8080",
	}

	for i := 0; i < 50; i++ {
		for i, k := range urls {
			bb.mq.RabbitMqAddMessages(JobConfig{
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
