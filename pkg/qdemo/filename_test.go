package qdemo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
)

func TestFilename_Mode(t *testing.T) {
	assert.Equal(t, "", qdemo.Filename("").Mode())
	assert.Equal(t, "2on2", qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Mode())
	assert.Equal(t, "duel", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Mode())
	assert.Equal(t, "duel", qdemo.Filename("duel_midair_holy_vs_si7h[bravado]261022-2255.mvd").Mode())
	assert.Equal(t, "ffa", qdemo.Filename("ffa_3[dm6]261022-2255_01.mvd").Mode())
}

func TestFilename_Submode(t *testing.T) {
	assert.Equal(t, "", qdemo.Filename("").Submode())
	assert.Equal(t, "", qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Submode())
	assert.Equal(t, "", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Submode())
	assert.Equal(t, "midair", qdemo.Filename("duel_midair_holy_vs_si7h[endif]261022-2255.mvd").Submode())
	assert.Equal(t, "midair", qdemo.Filename("2on2_midair_red_vs_blue[endif]261022-2255.mvd").Submode())
	assert.Equal(t, "", qdemo.Filename("ffa_3[dm6]261022-2255_01.mvd").Submode())
}

func TestFilename_Participants(t *testing.T) {
	assert.Equal(t, []string{}, qdemo.Filename("").Participants())
	assert.Equal(t, []string{"]a[", "]sr["}, qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Participants())
	assert.Equal(t, []string{"holy", "si7h"}, qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Participants())
	assert.Equal(t, []string{"holy", "si7h"}, qdemo.Filename("duel_midair_holy_vs_si7h[bravado]261022-2255.mvd").Participants())
	assert.Equal(t, []string{"3"}, qdemo.Filename("ffa_3[dm6]261022-2255_01.mvd").Participants())
}

func TestFilename_Map(t *testing.T) {
	assert.Equal(t, "", qdemo.Filename("").Map())
	assert.Equal(t, "", qdemo.Filename("duel_holy_vs_si7h[").Map())
	assert.Equal(t, "", qdemo.Filename("duel_holy_vs_si7h[]261022-2255.mvd").Map())
	assert.Equal(t, "dm4", qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Map())
	assert.Equal(t, "bravado", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Map())
	assert.Equal(t, "bravado", qdemo.Filename("duel_midair_holy_vs_si7h[bravado]261022-2255.mvd").Map())
	assert.Equal(t, "dm6", qdemo.Filename("ffa_3[dm6]261022-2255_01.mvd").Map())
}

func TestFilename_DateTime(t *testing.T) {
	assert.Equal(t, "", qdemo.Filename("").DateTime())
	assert.Equal(t, "231022-2126", qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").DateTime())
	assert.Equal(t, "261022-2255", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").DateTime())
	assert.Equal(t, "261022-2255", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.qwd").DateTime())
	assert.Equal(t, "261022-2255", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255_01.mvd").DateTime())
}

func TestFilename_Date(t *testing.T) {
	assert.Equal(t, "", qdemo.Filename("").Date())
	assert.Equal(t, "231022", qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Date())
	assert.Equal(t, "261022", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Date())
	assert.Equal(t, "261022", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.qwd").Date())
	assert.Equal(t, "261022", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255_01.mvd").Date())
}

func TestFilename_Time(t *testing.T) {
	assert.Equal(t, "", qdemo.Filename("").Time())
	assert.Equal(t, "2126", qdemo.Filename("2on2_]a[_vs_]sr[[dm4]231022-2126.mvd").Time())
	assert.Equal(t, "2255", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.mvd").Time())
	assert.Equal(t, "2255", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255.qwd").Time())
	assert.Equal(t, "2255", qdemo.Filename("duel_holy_vs_si7h[bravado]261022-2255_01.mvd").Time())
}

func TestFilename_ParseDateTime(t *testing.T) {
	timeByDateStr := func(dateStr string) time.Time {
		t, _ := time.Parse("20060102-1504", dateStr)
		return t
	}

	assert.Equal(t, time.Time{}, qdemo.Filename("").ParseDateTime("ymd"))
	assert.Equal(t, timeByDateStr("20010203-2255"), qdemo.Filename("duel_holy_vs_si7h[bravado]010203-2255.mvd").ParseDateTime("ymd"))
	assert.Equal(t, timeByDateStr("20030201-2255"), qdemo.Filename("duel_holy_vs_si7h[bravado]010203-2255.mvd").ParseDateTime("dmy"))
	assert.Equal(t, timeByDateStr("20010203-2255"), qdemo.Filename("duel_holy_vs_si7h[bravado]20010203-2255.mvd").ParseDateTime("Ymd"))
}
