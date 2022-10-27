package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type QtvServer struct {
	Address        string
	DemoDateFormat string
}

func (s *QtvServer) DemoStreamUrl(filename string) string {
	return fmt.Sprintf("http://%s/watch.qtv?demo=%s", s.Address, filename)
}

func (s *QtvServer) DemoDownloadUrl(filename string) string {
	return fmt.Sprintf("http://%s/dl/demos/%s", s.Address, filename)
}

type QtvHostedDemo struct {
	Filename DemoFilename
	Server   QtvServer
}

func (d *QtvHostedDemo) MarshalJSON() ([]byte, error) {
	type QtvHostedDemoExport struct {
		Filename      string `json:"filename"`
		ServerAddress string `json:"server_address"`
		StreamUrl     string `json:"stream_url"`
		DownloadUrl   string `json:"download_url"`
	}

	return json.Marshal(QtvHostedDemoExport{
		Filename:      string(d.Filename),
		ServerAddress: d.Server.Address,
		StreamUrl:     d.StreamUrl(),
		DownloadUrl:   d.DownloadUrl(),
	})
}

func (d *QtvHostedDemo) StreamUrl() string {
	return d.Server.DemoStreamUrl(string(d.Filename))
}

func (d *QtvHostedDemo) DownloadUrl() string {
	return d.Server.DemoDownloadUrl(string(d.Filename))
}

type DemoFilename string

func (d DemoFilename) Mode() string {
	strVal := string(d)

	indexFirstUnderScore := strings.IndexRune(strVal, '_')
	if -1 == indexFirstUnderScore {
		return ""
	}

	return strVal[0:indexFirstUnderScore]
}

func (d DemoFilename) Participants() []string {
	strVal := string(d)

	indexFirstUnderScore := strings.IndexRune(strVal, '_')
	if -1 == indexFirstUnderScore {
		return make([]string, 0)
	}

	indexOpenBracket := strings.LastIndexByte(strVal, '[')
	if -1 == indexOpenBracket {
		return make([]string, 0)
	}

	participantStr := strVal[indexFirstUnderScore+1 : indexOpenBracket]

	const vsNeedle = "_vs_"
	if strings.Contains(participantStr, vsNeedle) {
		return strings.SplitN(participantStr, vsNeedle, 2)
	}

	return []string{participantStr}
}

func (d DemoFilename) Map() string {
	strVal := string(d)

	indexOpenBracket := strings.LastIndexByte(strVal, '[')
	if -1 == indexOpenBracket {
		return ""
	}

	indexCloseBracket := strings.LastIndexByte(strVal, ']')
	if -1 == indexCloseBracket {
		return ""
	}

	if indexCloseBracket-indexOpenBracket < 1 {
		return ""
	}

	return strVal[indexOpenBracket+1 : indexCloseBracket]
}

func (d DemoFilename) DateTime() string {
	strVal := string(d)

	indexCloseBracket := strings.LastIndexByte(strVal, ']')
	if -1 == indexCloseBracket {
		return ""
	}

	return strVal[indexCloseBracket+1 : len(strVal)-len(".mvd")]
}

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
