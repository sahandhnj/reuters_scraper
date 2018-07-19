package main

import (
	. "git.sahand.cloud/sahand/reuters_scraper/model"
	"time"
	"fmt"
)

func findArticleLinks () []Article {
	var articles []Article
	var archives []Archive

	archives = findArchivesPerDay(2018,7,2)

	chArticles := make(chan Article)
	chFinished := make(chan bool) 

	for _, dayArchive := range archives {
		go dayArchive.Crawl(chArticles, chFinished)
	}

	for c := 0; c < len(archives); {
		select {
		case article := <-chArticles:
			articles= append(articles,article)
		case <-chFinished:
			c++
		}
	}

	close(chArticles)

	return articles
}

func findArchivesPerDay(year int, month time.Month, day int) (archive []Archive){
	date := time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
	now := time.Now()
	
	for !date.After(now.AddDate(0,0,-2)) {
		date = date.AddDate(0,0,1)
		archive = append(archive, Archive{"https://www.reuters.com/resources/archive/us/" + date.Format("20060102") + ".html", date})
	}

	return
}

func main() {
	articles := findArticleLinks()

	fmt.Println("\nFound", len(articles), "unique articles:")

	for _,article := range articles {
		article.ParseHTML()
		fmt.Printf("URL: %s \nTitle: %s\n %s \n\n\n\n", article.URL, article.Title , article.Content)
	}
}