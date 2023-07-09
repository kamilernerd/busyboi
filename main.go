package main

import (
	"context"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Crawlerbox struct {
	queueMsgs chan amqp.Delivery
}

func main() {
	bb := &Crawlerbox{
		queueMsgs: make(chan amqp.Delivery),
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
		go worker(m.Body, ctx)
	}
}
