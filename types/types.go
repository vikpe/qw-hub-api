package types

import (
	"time"
)

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

type NewsItem struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Url   string `json:"url"`
}

type TwitchStream struct {
	Channel       string    `json:"channel"`
	Url           string    `json:"url"`
	Title         string    `json:"title"`
	ViewerCount   int       `json:"viewers"`
	Language      string    `json:"language"`
	ClientName    string    `json:"client_name"`
	ServerAddress string    `json:"server_address"`
	StartedAt     time.Time `json:"started_at"`
}
