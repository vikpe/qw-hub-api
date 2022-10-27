package types

import (
	"strings"
)

type DemoFilename string

func (f DemoFilename) Mode() string {
	filename := string(f)

	// special modes
	if strings.HasPrefix(filename, "duel_midair") {
		return "duel_midair"
	}

	indexFirstUnderScore := strings.IndexRune(filename, '_')
	if -1 == indexFirstUnderScore {
		return ""
	}

	return filename[0:indexFirstUnderScore]
}

func (f DemoFilename) Participants() []string {
	filename := string(f)
	indexFrom := len(f.Mode()) + 1
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

func (f DemoFilename) Map() string {
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

func (f DemoFilename) DateTime() string {
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

func (f DemoFilename) Date() string {
	dateTime := f.DateTime()
	if -1 == strings.IndexRune(dateTime, '-') {
		return ""
	}

	return strings.SplitN(dateTime, "-", 2)[0]
}

func (f DemoFilename) Time() string {
	dateTime := f.DateTime()
	if -1 == strings.IndexRune(dateTime, '-') {
		return ""
	}

	return strings.SplitN(dateTime, "-", 2)[1]
}

func (f DemoFilename) Number() string {
	return "5"
}
