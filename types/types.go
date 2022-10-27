package types

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
