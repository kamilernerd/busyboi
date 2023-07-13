# Crawlerbox
Light, fast and scalable web scraper for structured data.

## Features
 - Structured data
 - Reads scraping jobs from rabbitmq
 - Respects robotx.txt
 - Renders javascript
 - Regex matching option
 - Re-schedules jobs when url is unreachable
 - Allowes grouping of elements Parent -> child

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
