package qtvscraper_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

func TestDemoFilenamesDirectoryIndexFallback(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://foo:28000/demo_filenames.txt", httpmock.NewStringResponder(404, "missing"))
	httpmock.RegisterResponder("GET", "http://foo:28000/demos/", httpmock.NewStringResponder(http.StatusOK, `<a href="duel_a_vs_b[dm4]20261016-0034.mvd">demo</a>`))
	httpmock.RegisterResponder("GET", "http://foo:28000/demolist/", httpmock.NewStringResponder(http.StatusOK, `<table id="demos"><tr><td class="name">duel_a_vs_b[dm4]20261016-0034.mvd</td></tr></table>`))
	server := qtvscraper.Server{Address: "foo:28000"}
	filenames, err := server.DemoFilenames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"duel_a_vs_b[dm4]20261016-0034.mvd"}, filenames)

	httpmock.RegisterResponder("GET", "http://foo:28000/demolist/", httpmock.NewStringResponder(404, "missing"))
	filenames, err = server.DemoFilenames()
	assert.Error(t, err)
	assert.Empty(t, filenames)
}
