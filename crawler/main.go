package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

type Message struct {
	ID     int
	Index  int
	Link   string
	Avatar string
	Author string
	Title  string
	Date   string
	HTML   string
	Text   string
}

var (
	postIndex          = 0
	postIds            = make(map[int]bool)
	amqpBroker         broker.Broker
	topicIndex         = 0
	topicIds           = make(map[int]bool)
	rxOk               = regexp.MustCompile(`(http|https)://x7bwsmcore5fmx56.onion`)
	errEnqueueRedirect = errors.New("redirection not followed")
)

const (
	amqpAddr = "amqp://localhost"
	topic    = "topic.crawler"
	proxy    = "http://localhost:8123"
	www      = "http://x7bwsmcore5fmx56.onion"
)

type Extender struct {
	gocrawl.DefaultExtender
}

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

				messageJSON, _ := json.Marshal(message)

				msg := &broker.Message{
					Header: map[string]string{
						"Index": fmt.Sprintf("%d", postIndex),
					},
					Body: messageJSON,
				}

				if err := amqpBroker.Publish(topic, msg); err != nil {
					log.Printf("[pub] failed: %v", err)
				} else {
					fmt.Printf("[pub] pubbed message #%d (index: %d)\n", message.ID, message.Index)
				}
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

// Run Crawler service
func Run() {
	// Initiate RabbitMQ broker
	amqpBroker = rabbitmq.NewBroker(
		broker.Addrs(amqpAddr),
	)

	if err := amqpBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := amqpBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	fmt.Println("Crawler service is running")

	// Run crawler
	opts := gocrawl.NewOptions(new(Extender))
	opts.CrawlDelay = 1 * time.Second
	opts.SameHostOnly = true
	// opts.RobotUserAgent = "APIs-Google"
	// opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"
	// opts.LogFlags = gocrawl.LogAll
	opts.MaxVisits = 100

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(www)

	select {}
}
