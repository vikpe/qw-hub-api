package qtvscraper_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

func TestServer_DemoDownloadUrl(t *testing.T) {
	server := qtvscraper.Server{Address: "foo:28000"}
	assert.Equal(t, "http://foo:28000/dl/demos/foo.mvd", server.DemoDownloadUrl("foo.mvd"))
}

func TestServer_DemoQtvplayUrl(t *testing.T) {
	server := qtvscraper.Server{Address: "foo:28000"}
	assert.Equal(t, "file:foo.mvd@foo:28000", server.DemoQtvplayUrl("foo.mvd"))
}

func TestServer_DemoFilenames(t *testing.T) {
	t.Run("dns fail", func(t *testing.T) {
		server := qtvscraper.Server{Address: "foo:28000", DemoDateFormat: "ymd"}
		filenames, err := server.DemoFilenames()
		assert.Empty(t, filenames)
		assert.ErrorContains(t, err, "failure in name resolution")
	})

	t.Run("http request fail", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://foo:28000/demos/", httpmock.NewStringResponder(404, "page not found"))

		server := qtvscraper.Server{Address: "foo:28000", DemoDateFormat: "ymd"}
		filenames, err := server.DemoFilenames()
		assert.Empty(t, filenames)
		assert.ErrorContains(t, err, "url not found")
	})

	t.Run("success", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockedRepsonseBody, _ := ioutil.ReadFile("./test_files/demos_1.html")
		responder := httpmock.NewBytesResponder(http.StatusOK, mockedRepsonseBody)
		httpmock.RegisterResponder("GET", "http://foo:28000/demos/", responder)

		server := qtvscraper.Server{Address: "foo:28000", DemoDateFormat: "ymd"}
		expectedFilenames := []string{
			"duel_holy_vs_si7h[aerowalk]261022-2234.mvd",
			"duel_igggy_vs_rasta[aerowalk]261022-2224.mvd",
			"4on4_blue_vs_red[dm3]261022-2206.mvd",
		}
		filenames, err := server.DemoFilenames()
		assert.Equal(t, expectedFilenames, filenames)
		assert.Nil(t, err)
	})
}

func TestScraper_Demos(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responseFoo, _ := ioutil.ReadFile("./test_files/demos_1.html")
	responderFoo := httpmock.NewBytesResponder(http.StatusOK, responseFoo)
	httpmock.RegisterResponder("GET", "http://foo:28000/demos/", responderFoo)

	responseBar, _ := ioutil.ReadFile("./test_files/demos_2.html")
	responderBar := httpmock.NewBytesResponder(http.StatusOK, responseBar)
	httpmock.RegisterResponder("GET", "http://bar:28000/demos/", responderBar)

	scraper := qtvscraper.NewScraper([]qtvscraper.Server{
		{Address: "foo:28000", DemoDateFormat: "dmy"},
		{Address: "bar:28000", DemoDateFormat: "ymd"},
		{Address: "INVALID:28000", DemoDateFormat: "ymd"},
	})

	demos := scraper.Demos()
	assert.Len(t, demos, 5)
}
