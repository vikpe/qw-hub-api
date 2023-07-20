package qdemo

import (
	"fmt"
	"strings"
	"time"
)

type Filename string

func (f Filename) Mode() string {
	filename := string(f)

	indexFirstUnderScore := strings.IndexRune(filename, '_')
	if -1 == indexFirstUnderScore {
		return ""
	}

	return filename[0:indexFirstUnderScore]
}

func (f Filename) Submode() string {
	modes := []string{"duel", "2on2", "4on4", "3on3"}
	submodes := []string{"midair", "bf"}
	filename := string(f)

	for _, mode := range modes {
		for _, submode := range submodes {
			prefix := fmt.Sprintf("%s_%s", mode, submode)

			if strings.HasPrefix(filename, prefix) {
				return submode
			}
		}
	}

	return ""
}

func (f Filename) Participants() []string {
	var indexFrom int
	submode := f.Submode()
	filename := string(f)

	if len(submode) > 0 {
		indexFrom = len(strings.Join([]string{f.Mode(), submode}, "_")) + 1
	} else {
		indexFrom = len(f.Mode()) + 1
	}

	indexTo := strings.LastIndexByte(filename, '[')
	if -1 == indexTo {
		return make([]string, 0)
	}

	participantStr := filename[indexFrom:indexTo]

	const vsNeedle = "_vs_"
	if strings.Contains(participantStr, vsNeedle) {
		return strings.SplitN(participantStr, vsNeedle, 2)
	}

	return []string{participantStr}
}

func (f Filename) Map() string {
	filename := string(f)
	indexOpenBracket := strings.LastIndexByte(filename, '[')
	if -1 == indexOpenBracket {
		return ""
	}

	indexCloseBracket := strings.LastIndexByte(filename, ']')
	if -1 == indexCloseBracket {
		return ""
	}

	if indexCloseBracket-indexOpenBracket <= 1 {
		return ""
	}

	return filename[indexOpenBracket+1 : indexCloseBracket]
}

func (f Filename) DateTime() string {
	filename := string(f)
	indexCloseBracket := strings.LastIndexByte(filename, ']')
	if -1 == indexCloseBracket {
		return ""
	}

	indexFrom := indexCloseBracket + 1
	length := strings.IndexAny(filename[indexFrom:], "_.") // until _x or .ext
	indexTo := indexFrom + length
	return filename[indexFrom:indexTo]
}

func (f Filename) Date() string {
	dateTime := f.DateTime()
	if -1 == strings.IndexRune(dateTime, '-') {
		return ""
	}

	return strings.SplitN(dateTime, "-", 2)[0]
}

func (f Filename) Time() string {
	dateTime := f.DateTime()
	if -1 == strings.IndexRune(dateTime, '-') {
		return ""
	}

	return strings.SplitN(dateTime, "-", 2)[1]
}

func (f Filename) ParseDateTime(dateFormat string) time.Time {
	layoutDate := dateFormatToTimeLayout(dateFormat)
	layoutTime := "1504" // hhmm
	layout := fmt.Sprintf("%s-%s", layoutDate, layoutTime)
	demoTime, err := time.Parse(layout, f.DateTime())

	if err != nil {
		return time.Time{}
	}

	return demoTime
}

func dateFormatToTimeLayout(dateFormat string) string {
	const YMD = "060102"
	const YYYYMMDD = "20060102"
	const DMY = "020106"

	switch dateFormat {
	case "Ymd":
		return YYYYMMDD
	case "dmy":
		return DMY
	default:
		return YMD
	}
}
