package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
)

func TestStatsForPort(t *testing.T) {
	scores := []laststats.Entry{
		{Port: 28501, Demo: "old"},
		{Port: 28502, Demo: "other"},
		{Port: 28501, Demo: "new"},
	}
	assert.Equal(t, []laststats.Entry{scores[0], scores[2]}, statsForPort(scores, 28501, 10))
	assert.Equal(t, []laststats.Entry{scores[2]}, statsForPort(scores, 28501, 1))
	assert.Equal(t, []laststats.Entry{scores[1]}, statsForPort(scores, 28502, 10))
	assert.Equal(t, []laststats.Entry{}, statsForPort(scores, 28503, 10))
	assert.Equal(t, []laststats.Entry{}, statsForPort(nil, 28501, 10))
}
