package model

import (
	"time"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type Archive struct {
	Link string
	Date time.Time
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	
	return
}

func (a *Archive) Crawl(articles chan Article, chFinished chan bool) {
	fmt.Println("Crawling \"" + a.Date.String() + "\" for urls (" + a.Link + ")")
	resp, err := http.Get(a.Link)

	defer func() {
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + a.Link + "\"")
		return
	}

	b := resp.Body
	defer b.Close() 

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(t)
			if !ok {
				continue
			}

			isValidArticle := strings.Index(url, "http://www.reuters.com/article") == 0
			if isValidArticle {
				articles <- Article{url, a.Date , "" , ""}
			}
		}
	}
}