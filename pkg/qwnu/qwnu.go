package qwnu

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vikpe/qw-hub-api/pkg/scrape"
)

const qwnuURL = "https://www.quakeworld.nu"

type Event struct {
	Title   string `json:"title"`
	Status  string `json:"status"`
	Date    string `json:"date"`
	WikiUrl string `json:"wiki_url"`
	LogoUrl string `json:"logo_url"`
}

type ForumPost struct {
	Title  string `json:"title"`
	Forum  string `json:"forum"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Url    string `json:"url"`
}

type NewsPost struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Url   string `json:"url"`
}

func Events(limit int) ([]Event, error) {
	wikiOverviewUrl := fmt.Sprintf("%s/wiki/Overview", qwnuURL)
	doc, err := scrape.ReadDocument(wikiOverviewUrl)

	if err != nil {
		return make([]Event, 0), err
	}

	events := make([]Event, 0)
	statuses := []string{"upcoming", "ongoing", "completed"}

	const indexLogoCell = 0
	const indexLinkCell = 1
	const indexDateCell = 2

	for t := range statuses {
		doc.Find(fmt.Sprintf("#%s", statuses[t])).Find("tr").Each(func(i int, s *goquery.Selection) {
			if 0 == i || i >= limit { // skip heading and limit to x items
				return
			}

			cells := s.Children()

			linkElement := cells.Eq(indexLinkCell).Find("a").First()
			linkRelHref := linkElement.AttrOr("href", "#")
			logoRelUrl := cells.Eq(indexLogoCell).Find("img").First().AttrOr("src", "")
			logoRelUrl = strings.Replace(logoRelUrl, "21px", "32px", 1) // use 32px size

			event := Event{
				Title:   linkElement.AttrOr("title", "[parse fail]"),
				Status:  statuses[t],
				Date:    strings.TrimSpace(cells.Eq(indexDateCell).Text()),
				WikiUrl: fmt.Sprintf("%s%s", qwnuURL, linkRelHref),
				LogoUrl: fmt.Sprintf("%s%s", qwnuURL, logoRelUrl),
			}
			events = append(events, event)
		})
	}

	return events, nil
}

func ForumPosts(limit int) ([]ForumPost, error) {
	doc, err := scrape.ReadDocument(qwnuURL)

	if err != nil {
		return make([]ForumPost, 0), err
	}

	forumPosts := make([]ForumPost, 0)
	doc.Find("#frmForumActivity").Find("a").Each(func(i int, s *goquery.Selection) {
		if i >= limit { // limit to x items
			return
		}

		forumParts := strings.Split(s.Find(".link_recent_forum").Text(), " in ")
		forumPost := ForumPost{
			Title:  s.Find("b").Text(),
			Forum:  forumParts[1],
			Author: s.Find("div.link_recent_author").Text()[len("By "):],
			Date:   forumParts[0],
			Url:    fmt.Sprintf("%s%s", qwnuURL, s.AttrOr("href", "#")),
		}
		forumPosts = append(forumPosts, forumPost)
	})

	return forumPosts, nil

}

func NewsPosts(limit int) ([]NewsPost, error) {
	newsUrl := fmt.Sprintf("%s/feeds/news.php", qwnuURL)
	doc, err := scrape.ReadDocument(newsUrl)

	if err != nil {
		return make([]NewsPost, 0), err
	}

	newsPosts := make([]NewsPost, 0)
	doc.Find("item").Each(func(i int, s *goquery.Selection) {
		if i >= limit { // limit to x items
			return
		}

		pubDate := s.Find("pubDate").Text()
		newsPost := NewsPost{
			Title: s.Find("title").Text(),
			Date:  pubDate[:len(pubDate)-len(" hh:mm:ss +0000")],
			Url:   s.Find("guid").Text(),
		}
		newsPosts = append(newsPosts, newsPost)
	})

	return newsPosts, nil
}
