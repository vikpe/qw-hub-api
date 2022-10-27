package types_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/vikpe/qw-hub-api/types"
)

func TestQtvServer(t *testing.T) {
	server := types.QtvServer{Address: "qw.foppa.dk:28000"}
	assert.Equal(t, "http://qw.foppa.dk:28000/dl/demos/foo.mvd", server.DemoDownloadUrl("foo.mvd"))
	assert.Equal(t, "file:foo.mvd@qw.foppa.dk:28000", server.DemoQtvplayUrl("foo.mvd"))
}

func TestQtvHostedDemo(t *testing.T) {
	demo := types.QtvHostedDemo{
		Filename: types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd"),
		Server: types.QtvServer{
			Address:        "troopers.fi:28000",
			DemoDateFormat: "dmy",
		},
	}

	assert.Equal(t, "http://troopers.fi:28000/dl/demos/duel_holy_vs_si7h[bravado]261022-2255.mvd", demo.DownloadUrl())
	assert.Equal(t, "file:duel_holy_vs_si7h[bravado]261022-2255.mvd@troopers.fi:28000", demo.QtvplayUrl())
	expectedTime, _ := time.Parse("020106-1504", "261022-2255")
	assert.Equal(t, expectedTime, demo.Time())
}
