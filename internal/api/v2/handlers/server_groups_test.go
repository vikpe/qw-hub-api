package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestServerGroup_Ip(t *testing.T) {
	server := qserver.GenericServer{
		Address: "qw.foppa.dk:27015",
	}

	group := handlers.ServerGroup{
		Servers: []qserver.GenericServer{server},
	}

	assert.Equal(t, "qw.foppa.dk", group.Ip())
}

func TestServerGroup_Geo(t *testing.T) {
	server := qserver.GenericServer{
		Geo: geo.Location{
			CC: "SE",
		},
	}
	group := handlers.ServerGroup{
		Servers: []qserver.GenericServer{server},
	}
	assert.Equal(t, server.Geo, group.Geo())
}

func TestServerGroup_Title(t *testing.T) {
	t.Run("single server", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Quake.se KTX"}},
			},
		}
		assert.Equal(t, "Quake.se KTX", group.Title())
	})

	t.Run("common prefix", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Quake.se KTX #1"}},
				{Settings: qsettings.Settings{"hostname": "Quake.se KTX #2"}},
				{Settings: qsettings.Settings{"hostname": "Quake.se QTV"}},
			},
		}
		assert.Equal(t, "Quake.se", group.Title())
	})

	t.Run("common suffix", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "KTX #1 @ Quake.se"}},
				{Settings: qsettings.Settings{"hostname": "KTX #2 @ Quake.se"}},
				{Settings: qsettings.Settings{"hostname": "QTV @ Quake.se"}},
			},
		}
		assert.Equal(t, "@ Quake.se", group.Title())
	})

	t.Run("no prefix/suffix in common", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Alpha"}, Address: "quake.se:28501"},
				{Settings: qsettings.Settings{"hostname": "Beta"}},
				{Settings: qsettings.Settings{"hostname": "Gamma"}},
			},
		}
		assert.Equal(t, "quake.se", group.Title())
	})
}
