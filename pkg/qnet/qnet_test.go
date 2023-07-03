package qnet_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qnet"
	"testing"
)

func TestToIpHostPort(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		hostport, err := qnet.ToIpHostPort("_INVALID_")
		assert.Equal(t, "", hostport)
		assert.NotNil(t, err)
	})

	t.Run("invalid hostname", func(t *testing.T) {
		hostport, err := qnet.ToIpHostPort("_INVALID_:80")
		assert.Equal(t, "", hostport)
		assert.NotNil(t, err)
	})

	t.Run("valid hostport (ip)", func(t *testing.T) {
		hostport, err := qnet.ToIpHostPort("1.1.1.1:80")
		assert.Equal(t, "1.1.1.1:80", hostport)
		assert.Nil(t, err)
	})

	t.Run("valid hostport (hostname)", func(t *testing.T) {
		hostport, err := qnet.ToIpHostPort("one.one.one.one:80")
		assert.Contains(t, [2]string{"1.0.0.1:80", "1.1.1.1:80"}, hostport)
		assert.Nil(t, err)
	})
}
