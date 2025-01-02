package twitch

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/nicklaw5/helix/v2"
)

const quakeGameId = "7348"

type Stream struct {
	Id              string    `json:"id"`
	Channel         string    `json:"channel"`
	Url             string    `json:"url"`
	Title           string    `json:"title"`
	ViewerCount     int       `json:"viewers"`
	Language        string    `json:"language"`
	ClientName      string    `json:"client_name"`
	ServerAddress   string    `json:"server_address"`
	StartedAt       time.Time `json:"started_at"`
	DurationMinutes int       `json:"duration_minutes"`
	IsFeatured      bool      `json:"is_featured"`
	GameName        string    `json:"game_name"`
	ProfileImageUrl string    `json:"profile_image_url"`
}

type StreamerIndex map[string]string

func (s StreamerIndex) UserLogins() []string {
	result := make([]string, 0)

	for userLogin := range s {
		result = append(result, userLogin)
	}

	return result
}

func (s StreamerIndex) GetClientName(stream helix.Stream) string {
	if _, ok := s[stream.UserLogin]; ok {
		return s[stream.UserLogin]
	}
	return stream.UserName
}

type Scraper struct {
	client       *helix.Client
	streamers    StreamerIndex
	helixStreams []helix.Stream
	helixUsers   []helix.User
	shouldStop   bool
	interval     int
}

func NewScraper(clientID string, userAccessToken string, streamers StreamerIndex) (*Scraper, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		UserAccessToken: userAccessToken,
	})

	if err != nil {
		fmt.Println("twitch client", err.Error())
		return &Scraper{}, err
	}

	return &Scraper{
		streamers:    streamers,
		client:       client,
		interval:     60, // todo: decrease
		shouldStop:   false,
		helixStreams: make([]helix.Stream, 0),
		helixUsers:   make([]helix.User, 0),
	}, nil
}

func (scraper *Scraper) Streams() []Stream {
	result := make([]Stream, 0)
	featuredLogins := scraper.streamers.UserLogins()

	for _, stream := range scraper.helixStreams {
		if stream.GameID != quakeGameId {
			continue
		}

		elems := Stream{
			ClientName:      scraper.streamers.GetClientName(stream),
			Id:              stream.UserID,
			Channel:         stream.UserName,
			Language:        stream.Language,
			Title:           stream.Title,
			ViewerCount:     stream.ViewerCount,
			Url:             fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin),
			ServerAddress:   "",
			StartedAt:       stream.StartedAt,
			DurationMinutes: int(time.Since(stream.StartedAt).Minutes()),
			IsFeatured:      slices.Contains(featuredLogins, stream.UserLogin),
			GameName:        stream.GameName,
			ProfileImageUrl: "",
		}

		for _, user := range scraper.helixUsers {
			if user.ID == stream.UserID {
				elems.ProfileImageUrl = strings.Replace(user.ProfileImageURL, "300x300", "70x70", 1)
				break
			}
		}

		result = append(result, elems)
	}

	return result
}

func (scraper *Scraper) Start() {
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
				// get streams
				streamsRes, streamsErr := scraper.client.GetStreams(&helix.StreamsParams{
					First:   20,
					Type:    "live",
					GameIDs: []string{quakeGameId},
				})

				if len(streamsRes.ErrorMessage) > 0 {
					fmt.Println("error fetching twitch streams:", streamsRes.ErrorMessage)
					return
				}

				if streamsErr != nil {
					fmt.Println("error fetching twitch streams", streamsErr)
					return
				}

				scraper.helixStreams = streamsRes.Data.Streams

				// get user info
				userIDs := make([]string, 0)
				for _, stream := range scraper.helixStreams {
					userIDs = append(userIDs, stream.UserID)
				}

				usersRes, usersErr := scraper.client.GetUsers(&helix.UsersParams{
					IDs: userIDs,
				})

				if len(usersRes.ErrorMessage) > 0 {
					fmt.Println("error fetching twitch users:", usersRes.ErrorMessage)
					return
				}

				if usersErr != nil {
					fmt.Println("error fetching twitch users", streamsErr)
					return
				}

				scraper.helixUsers = usersRes.Data.Users
			}
		}()

		if tick == scraper.interval {
			tick = 0
		}
	}
}

func (scraper *Scraper) Stop() {
	scraper.shouldStop = true
}
