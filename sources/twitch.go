package sources

import (
	"fmt"
	"time"

	"github.com/nicklaw5/helix"
)

const quakeGameId = "7348"

type TwitchScraper struct {
	client     *helix.Client
	channels   []string
	Streams    []helix.Stream
	shouldStop bool
	interval   int
}

func NewTwitchScraper(clientID string, userAccessToken string, channels []string) (TwitchScraper, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		UserAccessToken: userAccessToken,
	})

	if err != nil {
		fmt.Println("twitch client", err.Error())
		return TwitchScraper{}, err
	}

	return TwitchScraper{
		channels:   channels,
		client:     client,
		interval:   30,
		shouldStop: false,
		Streams:    make([]helix.Stream, 0),
	}, nil
}

func (s *TwitchScraper) Start() {
	s.shouldStop = false

	go func() {
		ticker := time.NewTicker(time.Duration(1) * time.Second)
		tick := -1

		for ; true; <-ticker.C {
			if s.shouldStop {
				return
			}

			tick++

			go func() {
				currentTick := tick
				isTimeToUpdate := currentTick%s.interval == 0

				if isTimeToUpdate {
					response, err := s.client.GetStreams(&helix.StreamsParams{
						First:      10,
						GameIDs:    []string{quakeGameId},
						UserLogins: s.channels,
					})

					if err != nil {
						fmt.Println("error fetching twitch streams")
						return
					}

					s.Streams = response.Data.Streams
				}
			}()

			if tick == s.interval {
				tick = 0
			}
		}
	}()
}

func (s *TwitchScraper) Stop() {
	s.shouldStop = true
}
