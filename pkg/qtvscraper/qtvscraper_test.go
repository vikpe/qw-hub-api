package qtvscraper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

func TestServer(t *testing.T) {
	server := qtvscraper.Server{Address: "qw.foppa.dk:28000"}
	assert.Equal(t, "http://qw.foppa.dk:28000/dl/demos/foo.mvd", server.DemoDownloadUrl("foo.mvd"))
	assert.Equal(t, "file:foo.mvd@qw.foppa.dk:28000", server.DemoQtvplayUrl("foo.mvd"))
}
