package demoh_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers/demoh"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

var demos = []qtvscraper.Demo{
	{Filename: "duel_red_vs_blue[dm6]221028-0130.mvd", QtvAddress: "alpha:28000"},
	{Filename: "2on2_foo_vs_bar[dm6]221028-0110.mvd", QtvAddress: "alpha:28000"},
	{Filename: "4on4_red_vs_blue[dm3]221028-0100.mvd", QtvAddress: "beta:28000"},
}

func TestFilterByQtvAddress(t *testing.T) {
	expect := []qtvscraper.Demo{
		{Filename: "duel_red_vs_blue[dm6]221028-0130.mvd", QtvAddress: "alpha:28000"},
		{Filename: "2on2_foo_vs_bar[dm6]221028-0110.mvd", QtvAddress: "alpha:28000"},
	}
	assert.Equal(t, expect, demoh.FilterByQtvAddress(demos, "alpha:28000"))
	assert.Empty(t, demoh.FilterByQtvAddress(demos, "gamma:28000"))
}

func TestFilterByMode(t *testing.T) {
	expect := []qtvscraper.Demo{
		{Filename: "duel_red_vs_blue[dm6]221028-0130.mvd", QtvAddress: "alpha:28000"},
	}
	assert.Equal(t, expect, demoh.FilterByMode(demos, "duel"))
	assert.Empty(t, demoh.FilterByMode(demos, "3on3"))
}

func TestFilterByQuery(t *testing.T) {
	expect := []qtvscraper.Demo{
		{Filename: "duel_red_vs_blue[dm6]221028-0130.mvd", QtvAddress: "alpha:28000"},
		{Filename: "2on2_foo_vs_bar[dm6]221028-0110.mvd", QtvAddress: "alpha:28000"},
	}
	assert.Equal(t, expect, demoh.FilterByQuery(demos, "dm6"))
	assert.Empty(t, demoh.FilterByMode(demos, "__FOO__"))
}

func TestFilterByParams(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		params := new(demoh.DemoParams)
		assert.Len(t, demoh.FilterByParams(demos, params), len(demos))
	})

	t.Run("combination of params", func(t *testing.T) {
		demos = []qtvscraper.Demo{
			{Filename: "duel_red_vs_blue[dm6]221028-0130.mvd", QtvAddress: "alpha:28000"},
			{Filename: "2on2_red_vs_blue[dm4]221028-0120.mvd", QtvAddress: "alpha:28000"},
			{Filename: "duel_red_vs_yellow[dm4]221028-0110.mvd", QtvAddress: "beta:28000"},
			{Filename: "duel_red_vs_yellow[dm4]221028-0100.mvd", QtvAddress: "alpha:28000"},
			{Filename: "duel_red_vs_green[dm4]221028-0100.mvd", QtvAddress: "beta:28000"},
		}

		params := demoh.DemoParams{
			Mode:       "duel",
			Query:      "dm4 green",
			QtvAddress: "beta:28000",
			Limit:      2,
		}

		expect := []qtvscraper.Demo{
			{Filename: "duel_red_vs_green[dm4]221028-0100.mvd", QtvAddress: "beta:28000"},
		}
		assert.Equal(t, expect, demoh.FilterByParams(demos, &params))
	})

	t.Run("limit", func(t *testing.T) {
		params := new(demoh.DemoParams)
		params.Limit = 2
		assert.Len(t, demoh.FilterByParams(demos, params), 2)
	})
}

func TestSubstringMatch(t *testing.T) {
	assert.False(t, demoh.SubstringMatch("", ""))
	assert.False(t, demoh.SubstringMatch("", "alpha"))
	assert.False(t, demoh.SubstringMatch("alpha beta gamma", ""))
	assert.False(t, demoh.SubstringMatch("alpha beta gamma", "delta"))
	assert.False(t, demoh.SubstringMatch("alpha", "alpha beta gamma"))
	assert.True(t, demoh.SubstringMatch("alpha beta gamma", "alp"))
	assert.True(t, demoh.SubstringMatch("alpha beta gamma", "alpha"))
	assert.True(t, demoh.SubstringMatch("alpha beta gamma", "alp bet"))
}
