package qwnu

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"github.com/vikpe/qw-hub-api/pkg/htmlparse"
)

const qwnuURL = "https://www.quakeworld.nu"

var wikiURL = "https://www.quakeworld.nu/wiki"
var wikiOverviewUrl = fmt.Sprintf("%s/Overview", wikiURL)

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

type GameInSpotlight struct {
	Participants string              `json:"participants"`
	Description  string              `json:"description"`
	Stream       GameInSpotlightLink `json:"stream"`
	Event        GameInSpotlightLink `json:"event"`
	Date         string              `json:"date"`
}

type NewsPost struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Url   string `json:"url"`
}

type WikiArticle struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Date  string `json:"date"`
}

func Events(limit int) ([]Event, error) {
	doc, err := htmlparse.GetDocument(wikiOverviewUrl)

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
			if 0 == i || len(events) >= limit { // skip heading and limit to x items
				return
			}

			cells := s.Children()

			var title string
			var wikiUrl string

			linkElement := cells.Eq(indexLinkCell).Find("a").First()
			if linkElement.Length() > 0 {
				title = linkElement.AttrOr("title", "[parse fail]")
				wikiUrl = fmt.Sprintf("%s%s", qwnuURL, linkElement.AttrOr("href", "#"))
			} else {
				title = strings.TrimSpace(cells.Eq(indexLinkCell).Text())
				wikiUrl = ""
			}

			logoRelUrl := cells.Eq(indexLogoCell).Find("img").First().AttrOr("src", "")
			logoRelUrl = strings.Replace(logoRelUrl, "21px", "32px", 1) // use 32px size

			event := Event{
				Title:   title,
				Status:  statuses[t],
				Date:    strings.TrimSpace(cells.Eq(indexDateCell).Text()),
				WikiUrl: wikiUrl,
				LogoUrl: fmt.Sprintf("%s%s", qwnuURL, logoRelUrl),
			}
			events = append(events, event)
		})
	}

	return events, nil
}

func ForumPosts(limit int) ([]ForumPost, error) {
	doc, err := htmlparse.GetDocument(qwnuURL)

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

type GameInSpotlightLink struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func GamesInSpotlight() ([]GameInSpotlight, error) {
	doc, err := htmlparse.GetDocument(wikiOverviewUrl)

	if err != nil {
		return make([]GameInSpotlight, 0), err
	}

	gamesDiv := doc.Find(".GameInSpotlight")

	descriptions := make([]string, 0)
	gamesDiv.ChildrenFiltered("div").Each(func(i int, s *goquery.Selection) {
		descriptions = append(descriptions, cleanHtmlText(s.Text()))
	})

	games := make([]GameInSpotlight, 0)

	gamesDiv.Find("tbody").Each(func(i int, s *goquery.Selection) {
		rows := s.ChildrenFiltered("tr")
		secondRow := rows.Next()

		// stream
		stream := GameInSpotlightLink{}
		streamLink := secondRow.Find(".twitchlink a")
		if streamLink.Length() > 0 {
			stream.Title = streamLink.Text()
			stream.Url = streamLink.AttrOr("href", "#")
		}

		// event
		eventLink := secondRow.ChildrenFiltered("td").Last().Find("a")
		event := GameInSpotlightLink{}

		if eventLink.Length() > 0 {
			event.Title = eventLink.Text()
			event.Url = wikiLinkHref(eventLink.AttrOr("href", "#"))
		}

		// result
		game := GameInSpotlight{
			Participants: cleanHtmlText(rows.First().Text()),
			Stream:       stream,
			Event:        event,
			Date:         secondRow.Find(".datetime").Text(),
		}

		// description
		if len(descriptions) >= i {
			game.Description = descriptions[i]
		}

		games = append(games, game)
	})

	return games, nil
}

func NewsPosts(limit int) ([]NewsPost, error) {
	newsUrl := fmt.Sprintf("%s/feeds/news.php", qwnuURL)
	doc, err := htmlparse.GetDocument(newsUrl)

	if err != nil {
		return make([]NewsPost, 0), err
	}

	newsPosts := make([]NewsPost, 0)
	doc.Find("item").Each(func(i int, s *goquery.Selection) {
		if i >= limit { // limit to x items
			return
		}

		pubDate := s.Find("pubDate").Text()
		dateFormat := " hh:mm:ss +0000"
		date := ""

		if len(pubDate) > len(dateFormat) {
			date = pubDate[:len(pubDate)-len(dateFormat)]
		}

		newsPost := NewsPost{
			Title: s.Find("title").Text(),
			Date:  date,
			Url:   s.Find("guid").Text(),
		}
		newsPosts = append(newsPosts, newsPost)
	})

	return newsPosts, nil
}

func WikiRecentChanges(limit int) ([]WikiArticle, error) {
	feedUrl := fmt.Sprintf("%s/w/api.php?hidebots=1&hidepreviousrevisions=1&namespace=0&urlversion=2&days=100&limit=20&action=feedrecentchanges&feedformat=rss", qwnuURL)
	doc, err := htmlparse.GetDocument(feedUrl)

	if err != nil {
		return make([]WikiArticle, 0), err
	}

	articles := make([]WikiArticle, 0)
	titles := make([]string, 0)
	doc.Find("item").Each(func(i int, s *goquery.Selection) {
		if len(articles) >= limit { // limit to x items
			return
		}

		article := WikiArticle{
			Title: s.ChildrenFiltered("title").Text(),
			Url:   strings.Replace(s.ChildrenFiltered("comments").Text(), "Talk:", "", 1),
			Date:  s.ChildrenFiltered("pubDate").Text(),
		}

		if lo.Contains(titles, article.Title) {
			return
		}

		articles = append(articles, article)
		titles = append(titles, article.Title)
	})

	return articles, nil
}

func cleanHtmlText(htmlText string) string {
	result := regexp.MustCompile(`\s+`).ReplaceAllString(htmlText, " ")
	result = strings.ReplaceAll(result, "\u00a0", "")
	result = strings.TrimSpace(result)
	return result
}

func wikiLinkHref(href string) string {
	if len(href) == 0 {
		return href
	} else if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s%s", qwnuURL, href)
	} else {
		return href
	}
}
