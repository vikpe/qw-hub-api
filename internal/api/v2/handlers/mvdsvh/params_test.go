package mvdsvh_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers/mvdsvh"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"testing"
)

var serverAlpha = mvdsv.Mvdsv{Address: "alpha:28501", PlayerSlots: slots.New(4, 0)}
var serverBeta = mvdsv.Mvdsv{Address: "beta:28501", PlayerSlots: slots.New(4, 2)}
var serverGamma = mvdsv.Mvdsv{Address: "gamma:28501", PlayerSlots: slots.New(4, 4)}
var servers = []mvdsv.Mvdsv{serverAlpha, serverBeta, serverGamma}

func TestFilterByEmpty(t *testing.T) {
	t.Run("include", func(t *testing.T) {
		assert.Equal(t, servers, mvdsvh.FilterByEmpty(servers, "include"))
		assert.Equal(t, servers, mvdsvh.FilterByEmpty(servers, ""))
	})

	t.Run("exclude", func(t *testing.T) {
		expect := []mvdsv.Mvdsv{serverBeta, serverGamma}
		assert.Equal(t, expect, mvdsvh.FilterByEmpty(servers, "exclude"))
	})

	t.Run("only", func(t *testing.T) {
		expect := []mvdsv.Mvdsv{serverAlpha}
		assert.Equal(t, expect, mvdsvh.FilterByEmpty(servers, "only"))
	})
}

func TestFilterByHostname(t *testing.T) {
	assert.Equal(t, servers, mvdsvh.FilterByHostname(servers, ""))
	assert.Equal(t, []mvdsv.Mvdsv{serverAlpha}, mvdsvh.FilterByHostname(servers, "alpha"))
	assert.Equal(t, []mvdsv.Mvdsv{serverAlpha}, mvdsvh.FilterByHostname(servers, "alpha:28501"))

	t.Run("by matching parsed hostname", func(t *testing.T) {
		var serverDelta = mvdsv.Mvdsv{
			Address: "1.1.1.1:28501",
			Settings: qsettings.Settings{
				"hostname_parsed": "delta:28501",
			},
		}

		assert.Equal(t, []mvdsv.Mvdsv{serverDelta}, mvdsvh.FilterByHostname([]mvdsv.Mvdsv{serverAlpha, serverDelta}, "delta"))
	})
}

func TestFilterByParams(t *testing.T) {
	params := mvdsvh.NewMvdsvParams()
	params.Empty = "exclude"
	params.Limit = 1
	assert.Equal(t, []mvdsv.Mvdsv{serverBeta}, mvdsvh.FilterByParams(servers, params))
}
