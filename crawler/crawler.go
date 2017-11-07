package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	instance *gocrawl.Crawler
}

type Extender struct {
	gocrawl.DefaultExtender
}

var (
	postIndex          = 0
	postIds            = make(map[int]bool)
	topicIndex         = 0
	topicIds           = make(map[int]bool)
	rxOk               = regexp.MustCompile(`(http|https)://x7bwsmcore5fmx56.onion`)
	errEnqueueRedirect = errors.New("redirection not followed")
)

func (x *Extender) Start(seeds interface{}) interface{} {
	fmt.Println("> Crawling process started")
	return seeds
}

func (x *Extender) End(err error) {
	fmt.Println("> Crawling done")
	fmt.Println("  topics:", postIndex)
	fmt.Println("  posts:", topicIndex)
}

func CheckRedirect(req *http.Request, via []*http.Request) error {
	if isRobotsURL(req.URL) {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		if len(via) > 0 {
			req.Header.Set("User-Agent", via[0].Header.Get("User-Agent"))
		}
		return nil
	}

	return errEnqueueRedirect
}

func (x *Extender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	if isTopic(ctx.NormalizedURL().String()) {
		doc.Find(".post").Each(func(i int, item *goquery.Selection) {
			permalink := item.Find(".posthead .post-link a.permalink")
			link, _ := permalink.Attr("href")
			id, _ := strconv.Atoi(strings.Split(strings.Split(link, "pid=")[1], "#")[0])
			avatar, _ := item.Find(".postbody .post-author .author-ident .useravatar img").Attr("src")
			content := item.Find(".postbody .post-entry .entry-content")
			html, _ := content.Html()

			if _, ok := topicIds[id]; !ok && id != 0 {
				postIndex++
				topicIds[id] = true

				message := Message{
					ID:     id,
					Index:  postIndex,
					Link:   link,
					Avatar: avatar,
					Author: item.Find(".posthead .post-byline strong").Text(),
					Title:  item.Find(".postbody .post-entry .entry-title").Text(),
					Date:   parseDate(permalink.Text()),
					HTML:   html,
					Text:   content.Text(),
				}

				ch <- message
			}
		})
	} else if isForum(ctx.NormalizedURL().String()) {
		doc.Find(".main-item .item-subject a").Each(func(i int, item *goquery.Selection) {
			link, _ := item.Attr("href")
			id, _ := strconv.Atoi(strings.Split(strings.Split(link, "id=")[1], "#")[0])

			if _, ok := postIds[id]; !ok && id != 0 {
				postIds[id] = true
				topicIndex++
			}
		})
	}

	// Return nil and true - let gocrawl find the links
	return nil, true
}

func (x *Extender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}

func (x *Extender) Fetch(ctx *gocrawl.URLContext, userAgent string, headRequest bool) (*http.Response, error) {
	var reqType string

	// Prepare the request with the right user agent
	if headRequest {
		reqType = "HEAD"
	} else {
		reqType = "GET"
	}

	req, e := http.NewRequest(reqType, ctx.NormalizedURL().String(), nil)

	if e != nil {
		return nil, e
	}

	req.Header.Set("User-Agent", userAgent)

	proxyURL, _ := url.Parse(proxy)

	httpClient := &http.Client{
		CheckRedirect: CheckRedirect,
		Transport:     &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}

	return httpClient.Do(req)
}

func (c *Crawler) init() {
	opts := gocrawl.NewOptions(new(Extender))
	opts.CrawlDelay = 1 * time.Second
	opts.SameHostOnly = true
	// opts.RobotUserAgent = "APIs-Google"
	// opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"
	// opts.LogFlags = gocrawl.LogAll
	// opts.MaxVisits = 100

	c.instance = gocrawl.NewCrawlerWithOptions(opts)
}

func (c *Crawler) run() {
	c.instance.Run(www)
}
