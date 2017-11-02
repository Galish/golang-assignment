package indexer

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func parseTokens(r io.Reader) []string {
	tokens := []string{}
	// separators := []rune{' ', ')', '('}

	d := html.NewTokenizer(r)
	rg, _ := regexp.Compile("^[^a-zA-Z0-9]*|[^a-zA-Z0-9]*$")

	for {
		// token type
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			return tokens
		}

		if tokenType == html.TextToken {
			token := strings.TrimSpace(d.Token().Data)

			// s := "my string(qq bb)zz"
			// ss := split(s, separators)
			// token := split(d.Token().Data, separators)

			for _, t := range strings.Fields(token) {
				// for _, t := range split(token, separators) {
				tTrimmed := rg.ReplaceAllString(t, "")
				if len(tTrimmed) != 0 && !isValidURL(tTrimmed) {
					tokens = append(tokens, strings.ToLower(tTrimmed))
				}
			}
		}
	}
}

func split(s string, separators []rune) []string {
	f := func(r rune) bool {
		for _, s := range separators {
			if r == s {
				return true
			}
		}
		return false
	}

	return strings.FieldsFunc(s, f)
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)

	return err == nil
}

func getKey(value int) string {
	return fmt.Sprintf("id#%d", value)
}

func inter(slices [][]string) []string {
	items := make(map[string]bool)
	res := []string{}

	for i, slice := range slices {
		for _, key := range slice {
			for m, _slice := range slices {
				if i != m {
					ct := contains(_slice, key)

					if _, ok := items[key]; ok == false {
						items[key] = ct
					} else if items[key] == true && ct == false {
						items[key] = false
					}
				}
			}
		}
	}

	for key := range items {
		if items[key] == true {
			res = append(res, key)
		}
	}

	return res
}

func merge(slices [][]string) []string {
	res := []string{}

	for _, slice := range slices {
		for _, item := range slice {
			if contains(res, item) == false {
				res = append(res, item)
			}
		}
	}

	return res
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
