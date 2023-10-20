package qtvscraper_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
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
		fmt.Println(err)
		assert.True(
			t,
			strings.Contains(err.Error(), "failure in name resolution") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "server misbehaving"),
		)
	})

	t.Run("http request fail", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://foo:28000/demo_filenames.txt", httpmock.NewStringResponder(404, "page not found"))
		httpmock.RegisterResponder("GET", "http://foo:28000/demos/", httpmock.NewStringResponder(404, "page not found"))

		server := qtvscraper.Server{Address: "foo:28000", DemoDateFormat: "ymd"}
		filenames, err := server.DemoFilenames()
		assert.Empty(t, filenames)
		assert.ErrorContains(t, err, "url not found")
	})

	t.Run("success", func(t *testing.T) {
		t.Run("by reading /demo_filenames.txt", func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			mockedRepsonseBody, _ := os.ReadFile("./test_files/demo_filenames.txt")
			responder := httpmock.NewBytesResponder(http.StatusOK, mockedRepsonseBody)
			httpmock.RegisterResponder("GET", "http://foo:28000/demo_filenames.txt", responder)
			httpmock.RegisterResponder("GET", "http://foo:28000/demos/", httpmock.NewStringResponder(404, "page not found"))

			server := qtvscraper.Server{Address: "foo:28000", DemoDateFormat: "Ymd"}
			expectedFilenames := []string{
				"2on2_red_vs_blue[dm4]20230708-1645.mvd",
				"2on2_blue_vs_red[dm4]20230708-1625.mvd",
				"2on2_blue_vs_red[dm4]20230708-1611.mvd",
			}
			filenames, err := server.DemoFilenames()
			assert.Equal(t, expectedFilenames, filenames)
			assert.Nil(t, err)
		})

		t.Run("by parsing /demos/", func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			mockedRepsonseBody, _ := os.ReadFile("./test_files/demos_1.html")
			responder := httpmock.NewBytesResponder(http.StatusOK, mockedRepsonseBody)
			httpmock.RegisterResponder("GET", "http://foo:28000/demo_filenames.txt", httpmock.NewStringResponder(404, "page not found"))
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
	})
}

func TestScraper_Demos(t *testing.T) {
	t.Run("default params", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		serverAlphaResponse, _ := os.ReadFile("./test_files/demos_1.html")
		httpmock.RegisterResponder(
			"GET", "http://alpha:28000/demos/",
			httpmock.NewBytesResponder(http.StatusOK, serverAlphaResponse),
		)

		serverBetaResponse, _ := os.ReadFile("./test_files/demos_2.html")
		httpmock.RegisterResponder(
			"GET", "http://beta:28000/demos/",
			httpmock.NewBytesResponder(http.StatusOK, serverBetaResponse),
		)

		serverGammaResponse := errors.New("fail")
		httpmock.RegisterResponder(
			"GET", "http://gamma:28000/demos/",
			httpmock.NewErrorResponder(serverGammaResponse),
		)

		scraper := qtvscraper.NewScraper([]qtvscraper.Server{
			{Address: "alpha:28000", DemoDateFormat: "dmy"},
			{Address: "beta:28000", DemoDateFormat: "ymd"},
			{Address: "gamma:28000", DemoDateFormat: "ymd"},
		})

		// check result
		demos := scraper.Demos()
		assert.Len(t, demos, 4)

		expectedFirstDemoTime, _ := time.Parse("060102-1504", "221028-0355")
		expectedFirstDemo := qtvscraper.Demo{
			QtvAddress:  "beta:28000",
			Time:        expectedFirstDemoTime,
			Filename:    "duel_alpha_vs_beta[dm6]221028-0355.mvd",
			DownloadUrl: "http://beta:28000/dl/demos/duel_alpha_vs_beta[dm6]221028-0355.mvd",
			QtvplayUrl:  "file:duel_alpha_vs_beta[dm6]221028-0355.mvd@beta:28000",
		}

		assert.Equal(t, expectedFirstDemo, demos[0])

		// check requests/cache
		assert.Equal(t, 3, httpmock.GetTotalCallCount())

		scraper.Demos() // should use cache (no new request)
		assert.Equal(t, 3, httpmock.GetTotalCallCount())

		scraper.CacheDuration = 0
		scraper.Demos() // should scrape again (2 new requests)
		assert.Equal(t, 6, httpmock.GetTotalCallCount())
	})

	t.Run("demo max age", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		const day = 24 * time.Hour
		const timeLayout = "060102-1504"
		demo1Timestamp := time.Now().Add(-5 * day).Format(timeLayout)
		demo2Timestamp := time.Now().Add(-25 * day).Format(timeLayout)

		serverResponseBody := fmt.Sprintf(`
<table id="demos" cellspacing="0">
      <thead>
        <tr>
          <th class="stream">stream</th>
          <th class="save">Download</th>
          <th class="name">Demoname</th>
          <th class="size">Size</th>
        </tr>
      </thead>
      <tbody>
        <tr class="even">
          <td class="name">duel_foo_vs_bar[aerowalk]%s.mvd</td><td class="size">1144 kB</td>
        </tr>
        <tr class="odd">
          <td class="name">duel_foo_vs_bar[aerowalk]%s.mvd</td><td class="size">1178 kB</td>
        </tr>
      </tbody>
    </table>
`, demo1Timestamp, demo2Timestamp)

		httpmock.RegisterResponder(
			"GET", "http://delta:28000/demos/",
			httpmock.NewStringResponder(http.StatusOK, serverResponseBody),
		)

		fmt.Println("serverResponseBody", serverResponseBody)

		scraper := qtvscraper.NewScraper([]qtvscraper.Server{
			{Address: "delta:28000", DemoDateFormat: "ymd"},
		})
		scraper.CacheDuration = 1 * time.Nanosecond

		// check result
		scraper.DemoMaxAge = 30 * day
		assert.Len(t, scraper.Demos(), 2)

		scraper.DemoMaxAge = 6 * day
		assert.Len(t, scraper.Demos(), 1)

		scraper.DemoMaxAge = 3 * day
		assert.Len(t, scraper.Demos(), 0)
	})
}

func TestIsRelevantDemo(t *testing.T) {
	testCases := map[string]bool{
		"duel_testcfg_vs_mj23[dm2].mvd":            false,
		"ffa_1[dm3]220101-2055.mvd":                false,
		"duel_foo_vs_bar[povdmm4]220101-2055.mvd":  false,
		"2on2_foo_vs_bar[povdmm4]220101-2055.mvd":  false,
		"2on2_foo_vs_bar[foo_dmm4]220101-2055.mvd": false,
		"2on2_foo_vs_bar[dmm4_foo]220101-2055.mvd": false,
		"duel_foo_vs_bro[dm4]221101-0445.mvd":      false,
		"duel_bro_vs_foo[dm4]221101-0445.mvd":      false,
		"duel_timber_vs_foo[dm4]221101-0445.mvd":   false,
		"duel_foo_vs_bar[endif]220101-2055.mvd":    false,

		"duel_foo_vs_bar[bravado]220101-2055.mvd": true,
		"2on2_blue_vs_red[dm3]220101-2055.mvd":    true,
		"ctf_blue_vs_red[ctf4]220101-2055.mvd":    true,
		"wipeout_blue_vs_red[dm3]220101-2055.mvd": true,
		"4on4_blue_vs_red[dm3]220101-2055.mvd":    true,
		"4on4_foo_vs_bar[povdmm4]220101-2055.mvd": true,
	}

	for filename, expect := range testCases {
		t.Run(filename, func(t *testing.T) {
			demoFilename := qdemo.Filename(filename)
			assert.Equal(t, expect, qtvscraper.IsRelevantDemo(demoFilename))
		})
	}
}
