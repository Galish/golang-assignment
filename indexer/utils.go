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
