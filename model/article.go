package model

import (
	"time"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"bytes"
)

type Article struct {
	URL string
	Date time.Time
	Title string
	Content string
}

func isBody(t html.Token) (ok bool){
	if t.Data == "div" {
		for _, a := range t.Attr {
			if a.Key == "class" && strings.HasPrefix(a.Val, "body_")  {
				ok = true
			}
		}
	}

	return
}

func (article *Article) ParseHTML() {
	resp, err := http.Get(article.URL)

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + article.URL + "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)
	foundBoy := false
	var contentStrem bytes.Buffer
	
	for {
		tt := z.Next()

		switch tt {
			case html.ErrorToken:
				return
			case html.StartTagToken:
				t := z.Token()
				
				isTitle := ( t.Data == "title")
				if isTitle {
					tt = z.Next()
					t := z.Token()
					article.Title = strings.Trim(t.Data, " ")
				}

				if !foundBoy {
					foundBoy = isBody(t)
				}
				
				isP := ( t.Data == "p")
				if isP && foundBoy{
					tt = z.Next()
					
					if tt == html.TextToken {
						t := z.Token()
						if !strings.HasPrefix(t.Data,"*") {
							contentStrem.WriteString("\n" + t.Data)
						}
					}
				}
		}

		article.Content = contentStrem.String()
	}
}