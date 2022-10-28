package htmlparse

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetDocument(url string) (*goquery.Document, error) {
	// request page
	res, err := http.Get(url)
	if err != nil {
		return &goquery.Document{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err := errors.New(fmt.Sprintf("url not found: %s (%d)", url, res.StatusCode))
		return &goquery.Document{}, err
	}

	// load document
	return goquery.NewDocumentFromReader(res.Body)
}
