package qwnu_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func TestCleanHtmlText(t *testing.T) {
	text := "\t  Kombat   Duel\n                                                        5"
	assert.Equal(t, "Kombat Duel 5", qwnu.CleanHtmlText(text))
}

func TestGamesInSpotlight(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	wikiIndexHtml, _ := os.ReadFile("./test_files/wiki_index.html")
	response := httpmock.NewStringResponder(200, string(wikiIndexHtml))
	httpmock.RegisterResponder("GET", "https://www.quakeworld.nu/wiki/Overview", response)

	games, err := qwnu.GamesInSpotlight()
	assert.Len(t, games, 3)
	assert.Nil(t, err)

	assert.Equal(t, qwnu.GameInSpotlight{
		Participants: "Bernkaoch vs. bps",
		Description:  "WB Round 3 - 22:00 cest",
		Stream: qwnu.GameInSpotlightLink{
			Title: "badsebitv",
			Url:   "http://twitch.tv/badsebitv",
		},
		Event: qwnu.GameInSpotlightLink{
			Title: "Kombat Duel 5",
			Url:   "https://www.quakeworld.nu/wiki/Kombat_Duel_5",
		},
		Date: "29 April 2023 20:00",
	}, games[0])

	assert.Equal(t, qwnu.GameInSpotlight{
		Participants: "Commando (Milton, Diki, Xantom, Xterm, HENU, Martin) vs. Slackers (ParadokS, zero, niw, mazer, xunito)",
		Description:  "Semifinal, 20:00 cet",
		Stream: qwnu.GameInSpotlightLink{
			Title: "badsebitv",
			Url:   "http://twitch.tv/badsebitv",
		},
		Event: qwnu.GameInSpotlightLink{
			Title: "Qlan War Tournament 5",
			Url:   "https://www.quakeworld.nu/wiki/Qlan_War_Tournament_5",
		},
		Date: "5 Nov 2023 19:00",
	}, games[2])
}

func TestEvents(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	wikiIndexHtml, _ := os.ReadFile("./test_files/wiki_index.html")
	response := httpmock.NewStringResponder(200, string(wikiIndexHtml))
	httpmock.RegisterResponder("GET", "https://www.quakeworld.nu/wiki/Overview", response)

	events, err := qwnu.Events(2)
	assert.Len(t, events, 2)
	assert.Nil(t, err)

	assert.Equal(t, qwnu.Event{
		Title:   "QW LAN PL 2023",
		Status:  "upcoming",
		Date:    "08 Nov",
		WikiUrl: "",
		LogoUrl: "https://www.quakeworld.nu/w/images/thumb/b/b8/Dqer-icon.png/32px-Dqer-icon.png",
	}, events[0])

	assert.Equal(t, qwnu.Event{
		Title:   "TEC Elite Cup 2023",
		Status:  "upcoming",
		Date:    "29 Sep",
		WikiUrl: "https://www.quakeworld.nu/wiki/TEC_Elite_Cup_2023",
		LogoUrl: "https://www.quakeworld.nu/w/images/9/9f/Tbd-icon.png",
	}, events[1])
}

func TestNewsPosts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	newsFeedXml, _ := os.ReadFile("./test_files/feed_news.xml")
	response := httpmock.NewStringResponder(200, string(newsFeedXml))
	httpmock.RegisterResponder("GET", "https://www.quakeworld.nu/feeds/news.php", response)

	newsPosts, err := qwnu.NewsPosts(2)
	assert.Len(t, newsPosts, 2)
	assert.Nil(t, err)

	assert.Equal(t, qwnu.NewsPost{
		Title: "QHLAN 2024 - Signups Open",
		Date:  "Tue, 06 Jun 2023",
		Url:   "https://www.quakeworld.nu/news/1185",
	}, newsPosts[0])
}

func TestForumPosts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	qwnuIndexHtml, _ := os.ReadFile("./test_files/qwnu_index.html")
	response := httpmock.NewStringResponder(200, string(qwnuIndexHtml))
	httpmock.RegisterResponder("GET", "https://www.quakeworld.nu", response)

	forumPosts, err := qwnu.ForumPosts(2)
	assert.Len(t, forumPosts, 2)
	assert.Nil(t, err)

	assert.Equal(t, qwnu.ForumPost{
		Title:  "Map \"trick\", last stage",
		Forum:  "General Discussion",
		Author: "JSS",
		Date:   "2 days ago",
		Url:    "https://www.quakeworld.nu/forum/topic/7690/110144/map-trick-last-stage#110144",
	}, forumPosts[0])
}

func TestWikiRecentChanges(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	qwnuIndexHtml, _ := os.ReadFile("./test_files/feed_wiki_recent_changes.xml")
	response := httpmock.NewStringResponder(200, string(qwnuIndexHtml))
	httpmock.RegisterResponder("GET", "https://www.quakeworld.nu/w/api.php?hidebots=1&hidepreviousrevisions=1&namespace=0&urlversion=2&days=100&limit=20&action=feedrecentchanges&feedformat=rss", response)

	articles, err := qwnu.WikiRecentChanges(5)
	assert.Len(t, articles, 5)
	assert.Nil(t, err)

	assert.Equal(t, qwnu.WikiArticle{
		Title: "Brunowa",
		Url:   "https://www.quakeworld.nu/wiki/Brunowa",
		Date:  "Mon, 10 Jul 2023 05:42:05 GMT",
	}, articles[0])
}
