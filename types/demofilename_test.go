package types_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/vikpe/qw-hub-api/types"
)

func TestDemoFilename_Mode(t *testing.T) {
	assert.Equal(t, "", types.DemoFilename("").Mode())
	assert.Equal(t, "2on2", types.DemoFilename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Mode())
	assert.Equal(t, "duel", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Mode())
	assert.Equal(t, "duel_midair", types.DemoFilename("duel_midair_holy_vs_si7h[bravado]261022-2255.mvd").Mode())
	assert.Equal(t, "ffa", types.DemoFilename("ffa_3[dm6]261022-2255_01.mvd").Mode())
}

func TestDemoFilename_Participants(t *testing.T) {
	assert.Equal(t, []string{}, types.DemoFilename("").Participants())
	assert.Equal(t, []string{"]a[", "]sr["}, types.DemoFilename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Participants())
	assert.Equal(t, []string{"holy", "si7h"}, types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Participants())
	assert.Equal(t, []string{"holy", "si7h"}, types.DemoFilename("duel_midair_holy_vs_si7h[bravado]261022-2255.mvd").Participants())
	assert.Equal(t, []string{"3"}, types.DemoFilename("ffa_3[dm6]261022-2255_01.mvd").Participants())
}

func TestDemoFilename_Map(t *testing.T) {
	assert.Equal(t, "", types.DemoFilename("").Map())
	assert.Equal(t, "", types.DemoFilename("duel_holy_vs_si7h[").Map())
	assert.Equal(t, "", types.DemoFilename("duel_holy_vs_si7h[]261022-2255.mvd").Map())
	assert.Equal(t, "dm4", types.DemoFilename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Map())
	assert.Equal(t, "bravado", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Map())
	assert.Equal(t, "bravado", types.DemoFilename("duel_midair_holy_vs_si7h[bravado]261022-2255.mvd").Map())
	assert.Equal(t, "dm6", types.DemoFilename("ffa_3[dm6]261022-2255_01.mvd").Map())
}

func TestDemoFilename_DateTime(t *testing.T) {
	assert.Equal(t, "", types.DemoFilename("").DateTime())
	assert.Equal(t, "231022-2126", types.DemoFilename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").DateTime())
	assert.Equal(t, "261022-2255", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd").DateTime())
	assert.Equal(t, "261022-2255", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.qwd").DateTime())
	assert.Equal(t, "261022-2255", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255_01.mvd").DateTime())
}

func TestDemoFilename_Date(t *testing.T) {
	assert.Equal(t, "", types.DemoFilename("").Date())
	assert.Equal(t, "231022", types.DemoFilename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Date())
	assert.Equal(t, "261022", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Date())
	assert.Equal(t, "261022", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.qwd").Date())
	assert.Equal(t, "261022", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255_01.mvd").Date())
}

func TestDemoFilename_Time(t *testing.T) {
	assert.Equal(t, "", types.DemoFilename("").Time())
	assert.Equal(t, "2126", types.DemoFilename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Time())
	assert.Equal(t, "2255", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Time())
	assert.Equal(t, "2255", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255.qwd").Time())
	assert.Equal(t, "2255", types.DemoFilename("duel_holy_vs_si7h[bravado]261022-2255_01.mvd").Time())
}
