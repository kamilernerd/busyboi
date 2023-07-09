package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
 * DOM parser
 * Parses HTML downloaded by the worker and returns elements
 * specified in a json config.
 */
func parse(conf []JobConfigField, DOM string) map[string]interface{} {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(DOM))

	if err != nil {
		return nil
	}

	var fields = make(map[string]interface{}, len(conf))
	for _, field := range conf {
		if field.Attr != "" {
			attr, ok := doc.Find(field.Selector).Attr(field.Attr)
			if ok {
				fields[field.Name] = attr
			}
		} else {
			fields[field.Name] = doc.Find(field.Selector).Text()
		}

		if field.Children != nil {
			fields[field.Name] = parse(field.Children, DOM)
		}

		if field.Filters != nil {
			for _, filter := range field.Filters {
				switch filter {
				case "number":
				case "trim":
					break
				}
			}
		}
	}

	return fields
}
