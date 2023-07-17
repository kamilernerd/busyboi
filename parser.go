package main

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
 * DOM parser
 * Parses HTML downloaded by the worker and returns elements
 * specified in a json config.
 *
 * Regex match string when provided in a "parent" overrides
 * the regex string provided in the children elements
 */
func parse(conf []JobConfigField, DOM string, parentRegex string) map[string]interface{} {
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
			fields[field.Name] = parse(field.Children, DOM, field.Regex)
			continue
		}

		if parentRegex != "" {
			field.Regex = parentRegex
		}

		if field.Regex != "" {
			c, err := regexp.Compile(field.Regex)
			if err != nil {
				// Log error
			}

			f := c.FindString(fields[field.Name].(string))

			if f == "" {
				// Log error
			}

			fields[field.Name] = f
		}
	}

	return fields
}
