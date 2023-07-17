# Busyboi
Light, fast and scalable web scraper for structured data.

## Features
 - Structured data
 - Reads scraping jobs from rabbitmq
 - Respects robotx.txt
 - Renders javascript
 - Regex matching option
 - Re-schedules jobs when url is unreachable
 - Allowes grouping of elements Parent -> child
 - Config based
 - Workload limiter (variable)

## How to install?
 ```
 $ git clone git@github.com:kamilernerd/busyboi.git
 $ cd busyboi
 $ make
 ```
 
## CLI arguments
 ```
Usage of ./busyboi:
  -concurreny_limit int
    	Limits how many crawling jobs can run simultaneously (default 10)
  -queue_host string
    	Hostname for rabbitmq (default "localhost")
  -queue_name string
    	Queue name for rabbitmq (default "busyboi")
  -queue_password string
    	Password for rabbitmq (default "guest")
  -queue_port string
    	Port for rabbitmq (default "5672")
  -queue_user string
    	User for rabbitmq (default "guest")
 ```


## Example config
 ```
 {
    "collection": "some_random",
    "url": "somerandomurl.com/some/directory/index.html",
    "fields": [
        {
            "name": "some html element",
            "selector": "body > p[class=\"hello_world\"]"
        },
        {
            "name": "some html a element",
            "selector": "body > a[class=\"hello_world\"]",
            "attr": "href"
        },
        {
            "name": "a group of elements within a parent",
            "selector": ".some-parent-element",
            "children": [
                {
                    "name": "a link within a parent element",
                    "selector": "a[class=\"some_class"\]",
                    "attr": "href"
                }
            ]
        },
        {
            "name": "a group of elements within a parent",
            "selector": ".some-parent-element",
            "children": [
                {
                    "name": "a link within a parent element",
                    "selector": "a[class=\"some_class"\]",
                    "attr": "href",
                    "regex": "(http|s)+"
                }
            ]
        },
        {
            "name": "a group of elements within a parent",
            "selector": ".some-parent-element",
            "regex": "(http|s)+",
            "children": [
                {
                    "name": "a link within a parent element",
                    "selector": "a[class=\"some_class"\]",
                    "attr": "href"
                }
            ]
        }

    ]
}
```

## TODO
- In-memory TTL cache (DOM and parsed content)
- Maybe some long term storage support (MAYBE!)