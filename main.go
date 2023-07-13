package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExecutionsPerWorker = 50
)

type Crawlerbox struct {
	queueMsgs      chan amqp.Delivery
	wg             sync.WaitGroup
	workersCounter uint
}

func main() {
	bb := &Crawlerbox{
		queueMsgs:      make(chan amqp.Delivery),
		wg:             sync.WaitGroup{},
		workersCounter: 1,
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
			fmt.Printf("Worker #%s backing off for %s\n", m.MessageId, timeoutUntilNext.String())
			time.Sleep(timeoutUntilNext)
		}
		go worker(m.Body, ctx, bb)
	}

	bb.wg.Wait()
	close(bb.queueMsgs)
}
