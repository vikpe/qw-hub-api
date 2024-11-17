package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qversion"
)

func TestServerGroup_Host(t *testing.T) {
	server := qserver.GenericServer{
		Address: "qw.foppa.dk:27015",
	}

	group := handlers.ServerGroup{
		Servers: []qserver.GenericServer{server},
	}

	assert.Equal(t, "qw.foppa.dk", group.Host())
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

func TestServerGroup_Name(t *testing.T) {
	qtvServer := qserver.GenericServer{
		Settings: qsettings.Settings{"hostname": "QTV"},
		Version:  qversion.New("QTV"),
	}

	t.Run("no servers", func(t *testing.T) {
		group := handlers.ServerGroup{}
		assert.Equal(t, "", group.Name())
	})

	t.Run("single server", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Quake.se FFA"}},
			},
		}
		assert.Equal(t, "Quake.se FFA", group.Name())
	})

	t.Run("single game server", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Quake.se KTX"}},
			},
		}
		assert.Equal(t, "Quake.se KTX", group.Name())
	})

	t.Run("common prefix/suffix", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Quake.se KTX #1"}},
				{Settings: qsettings.Settings{"hostname": "Quake.se KTX #2"}},
				{Settings: qsettings.Settings{"hostname": "Quake.se FFA"}},
				qtvServer,
			},
		}
		assert.Equal(t, "Quake.se", group.Name())
	})

	t.Run("no prefix/suffix in common", func(t *testing.T) {
		group := handlers.ServerGroup{
			Servers: []qserver.GenericServer{
				{Settings: qsettings.Settings{"hostname": "Alpha"}, Address: "quake.se:28501"},
				{Settings: qsettings.Settings{"hostname": "Beta"}},
				qtvServer,
			},
		}
		assert.Equal(t, "quake.se", group.Name())
	})
}

func TestGetCommonHostname(t *testing.T) {
	t.Run("empty list of hostnames", func(t *testing.T) {
		hostnames := []string{}
		assert.Equal(t, "", handlers.GetCommonHostname(hostnames, 3))
	})

	t.Run("single hostname", func(t *testing.T) {
		hostnames := []string{"Quake.se KTX"}
		assert.Equal(t, "Quake.se KTX", handlers.GetCommonHostname(hostnames, 3))
	})

	t.Run("common prefix", func(t *testing.T) {
		hostnames := []string{
			"Quake.se KTX #1",
			"Quake.se KTX #2",
			"Quake.se FFA",
		}
		assert.Equal(t, "Quake.se", handlers.GetCommonHostname(hostnames, 3))
	})

	t.Run("common prefix (strip port)", func(t *testing.T) {
		hostnames := []string{
			"troopers.fi:28001",
			"troopers.fi:28002",
		}
		assert.Equal(t, "troopers.fi", handlers.GetCommonHostname(hostnames, 3))
	})

	t.Run("common suffix", func(t *testing.T) {
		hostnames := []string{
			"KTX #1 @ Quake.se",
			"KTX #2 @ Quake.se",
			"FFA @ Quake.se",
		}
		assert.Equal(t, "@ Quake.se", handlers.GetCommonHostname(hostnames, 3))
	})

	t.Run("no common prefix or suffix", func(t *testing.T) {
		hostnames := []string{
			"Alpha",
			"Beta",
			"Gamma",
		}
		assert.Equal(t, "", handlers.GetCommonHostname(hostnames, 3))
	})
}
