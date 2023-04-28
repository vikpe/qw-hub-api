package qwnu_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func TestGamesInSpotlight(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	wikiIndexHtml, _ := os.ReadFile("./test_files/wiki_index.html")
	response := httpmock.NewStringResponder(200, string(wikiIndexHtml))
	httpmock.RegisterResponder("GET", "https://www.quakeworld.nu/wiki/Overview", response)

	games, err := qwnu.GamesInSpotlight()
	assert.Len(t, games, 2)
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
}
