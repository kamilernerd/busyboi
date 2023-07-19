package main

import (
	"time"
)

type Bucket struct {
	Url       string
	Refreshed time.Time
	Data      map[string]interface{}
}

type Cache struct {
	ttl   time.Duration
	table map[string]Bucket
}

func NewCache(timeToLive time.Duration) *Cache {
	return &Cache{
		ttl:   timeToLive,
		table: make(map[string]Bucket),
	}
}

func (c *Cache) AddOrUpdate(bucket Bucket) {
	c.table[bucket.Url] = Bucket{}
	c.table[bucket.Url] = bucket
}

func (c *Cache) Remove(url string) {
	delete(c.table, url)
}

func (c *Cache) Has(url string) map[string]interface{} {
	if hit, ok := c.table[url]; ok {
		return hit.Data
	}
	return nil
}

func (c *Cache) Invalidator() {
	for {
		if len(c.table) == 0 {
			continue
		}

		func() {
			for _, v := range c.table {
				if v.Refreshed.Add(c.ttl).Compare(time.Now()) == -1 {
					c.Remove(v.Url)
				}
			}
		}()
	}
}
