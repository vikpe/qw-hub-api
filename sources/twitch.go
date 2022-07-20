package sources

import (
	"fmt"
	"time"

	"github.com/nicklaw5/helix"
)

type StreamerIndex map[string]string

func (s StreamerIndex) UserLogins() []string {
	result := make([]string, 0)

	for userLogin := range s {
		result = append(result, userLogin)
	}

	return result
}

type TwitchStream struct {
	Channel       string `json:"channel"`
	Url           string `json:"url"`
	Title         string `json:"title"`
	ViewerCount   int    `json:"viewers"`
	Language      string `json:"language"`
	ClientName    string `json:"client_name"`
	ServerAddress string `json:"server_address"`
}

type TwitchScraper struct {
	client       *helix.Client
	streamers    StreamerIndex
	helixStreams []helix.Stream
	shouldStop   bool
	interval     int
}

func (scraper TwitchScraper) Streams() []TwitchStream {
	result := make([]TwitchStream, 0)

	for _, stream := range scraper.helixStreams {
		result = append(result, TwitchStream{
			ClientName:    scraper.streamers[stream.UserLogin],
			Channel:       stream.UserName,
			Language:      stream.Language,
			Title:         stream.Title,
			ViewerCount:   stream.ViewerCount,
			Url:           fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin),
			ServerAddress: "",
		})
	}

	return result
}

func NewTwitchScraper(clientID string, userAccessToken string, streamers StreamerIndex) (*TwitchScraper, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		UserAccessToken: userAccessToken,
	})

	if err != nil {
		fmt.Println("twitch client", err.Error())
		return &TwitchScraper{}, err
	}

	return &TwitchScraper{
		streamers:    streamers,
		client:       client,
		interval:     5,
		shouldStop:   false,
		helixStreams: make([]helix.Stream, 0),
	}, nil
}

func (scraper *TwitchScraper) Start() {
	scraper.shouldStop = false

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	tick := -1

	for ; true; <-ticker.C {
		if scraper.shouldStop {
			return
		}

		tick++

		go func() {
			currentTick := tick
			isTimeToUpdate := currentTick%scraper.interval == 0

			if isTimeToUpdate {
				const quakeGameId = "7348"

				response, err := scraper.client.GetStreams(&helix.StreamsParams{
					First:      10,
					GameIDs:    []string{quakeGameId},
					UserLogins: scraper.streamers.UserLogins(),
				})

				if len(response.ErrorMessage) > 0 {
					fmt.Println("error fetching twitch streams:", response.ErrorMessage)
					return
				}

				if err != nil {
					fmt.Println("error fetching twitch streams", err)
					return
				}

				scraper.helixStreams = response.Data.Streams
			}
		}()

		if tick == scraper.interval {
			tick = 0
		}
	}
}

func (scraper *TwitchScraper) Stop() {
	scraper.shouldStop = true
}
