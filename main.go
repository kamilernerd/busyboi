package main

import (
	"context"
	"flag"
	"time"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Busyboi struct {
	queueMsgs chan amqp.Delivery
	mq        *Rabbitmq
	cache     *Cache
}

func main() {
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

	limiter := make(chan int, concurreny_limit)
	for m := range bb.queueMsgs {
		limiter <- 1

		go func(m amqp.Delivery) {
			worker(m.Body, ctx, bb)
			<-limiter
		}(m)
	}

	close(bb.queueMsgs)
}
