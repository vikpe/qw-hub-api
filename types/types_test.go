package types_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/vikpe/qw-hub-api/types"
)

func TestDemoFileName(t *testing.T) {
	t.Run("invalid filename", func(t *testing.T) {
		demoFileName := types.DemoFileName("")
		assert.Equal(t, "", demoFileName.Mode())
		assert.Equal(t, make([]string, 0), demoFileName.Participants())
		assert.Equal(t, "", demoFileName.Map())
		assert.Equal(t, "", demoFileName.DateTime())
	})

	t.Run("valid filename", func(t *testing.T) {
		t.Run("duel", func(t *testing.T) {
			demoFileName := types.DemoFileName("duel_holy_vs_si7h[bravado]261022-2255.mvd")
			assert.Equal(t, "duel", demoFileName.Mode())
			assert.Equal(t, []string{"holy", "si7h"}, demoFileName.Participants())
			assert.Equal(t, "bravado", demoFileName.Map())
			assert.Equal(t, "261022-2255", demoFileName.DateTime())
		})

		t.Run("ffa", func(t *testing.T) {
			demoFileName := types.DemoFileName("ffa_5[ztndm3]151022-2104.mvd")
			assert.Equal(t, "ffa", demoFileName.Mode())
			assert.Equal(t, []string{"5"}, demoFileName.Participants())
			assert.Equal(t, "ztndm3", demoFileName.Map())
			assert.Equal(t, "151022-2104", demoFileName.DateTime())
		})
	})
}

func TestQtvServer(t *testing.T) {
	server := types.QtvServer{Address: "qw.foppa.dk:28000"}
	assert.Equal(t, "http://qw.foppa.dk:28000/dl/demos/foo.mvd", server.DemoDownloadUrl("foo.mvd"))
	assert.Equal(t, "http://qw.foppa.dk:28000/watch.qtv?demo=foo.mvd", server.DemoStreamUrl("foo.mvd"))
}

func TestQtvHostedDemo(t *testing.T) {
	demo := types.QtvHostedDemo{
		Filename: types.DemoFileName("duel_holy_vs_si7h[bravado]261022-2255.mvd"),
		Server:   types.QtvServer{Address: "troopers.fi:28000"},
	}

	assert.Equal(t, "http://troopers.fi:28000/dl/demos/duel_holy_vs_si7h[bravado]261022-2255.mvd", demo.DownloadUrl())
	assert.Equal(t, "http://troopers.fi:28000/watch.qtv?demo=duel_holy_vs_si7h[bravado]261022-2255.mvd", demo.StreamUrl())
}
