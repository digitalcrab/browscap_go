package browscap_go

import (
	"fmt"
)

var (
	dict		*dictionary
	initialized bool
)

func InitBrowsCap(path string, force bool) error {
	if initialized && !force {
		return nil
	}
	var err	error

	// Load ini file
	if dict, err = loadFromIniFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	// Process ini file
	if err = dict.buildExpressions(); err != nil {
		return fmt.Errorf("browscap: An error occurred while processing data, %v ", err)
	}

	initialized = true
	return nil
}

func GetBrowser(userAgent string) (browser *Browser, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			browser = nil
			ok = false
		}
	}()

	for i, exp := range dict.expressions {
		if exp.Match([]byte(userAgent)) {
			sec := dict.sorted[i]
			data := dict.findData(sec.Name)
			browser = extractBrowser(data)
			ok = true
			return
		}
	}

	return
}

/*import (
	"strings"
	"regexp"
	"fmt"
)

const (
	DefaultPropertiesKey = "DefaultProperties"
)

var (
	parsed bool
	patterns Patterns
	defaultProperties = map[string]string{
		"Comment":				"DefaultProperties",
		"Browser":				"DefaultProperties",
		"Version":				"0.0",
		"MajorVer":				"0",
		"MinorVer":				"0",
		"Platform":				"",
		"Platform_Version":		"",
		"Alpha":				"false",
		"Beta":					"false",
		"Win16":				"false",
		"Win32":				"false",
		"Win64":				"false",
		"Frames":				"false",
		"IFrames":				"false",
		"Tables":				"false",
		"Cookies":				"false",
		"BackgroundSounds":		"false",
		"JavaScript":			"false",
		"VBScript":				"false",
		"JavaApplets":			"false",
		"ActiveXControls":		"false",
		"isMobileDevice":		"false",
		"isTablet":				"false",
		"isSyndicationReader":	"false",
		"Crawler":				"false",
		"CssVersion":			"0",
		"AolVersion":			"0",
	}


)

type Patterns []*Pattern

type Pattern struct {
	Pattern	string
	Data	map[string]string
	Regexp	*regexp.Regexp
}

type PatternsMap map[string]*Pattern

func InitBrowsCapXml(path string, force bool) error {
	if parsed && !force {
		return nil
	}

	var (
		caps	*xmlBrowserCaps
		err		error
		pMap	PatternsMap
	)

	if caps, err = parseXml(path); err != nil {
		return err
	}

	pMap = buildFromXml(caps)
	if patterns, err = rebase(pMap); err != nil {
		return fmt.Errorf("browscap: error on rebase, %v", err)
	}

	parsed = true

	return nil
}

func GetBrowser(agent string) (map[string]string, bool) {
	if parsed == true {
		for _, p := range patterns {
			if p.Regexp == nil {
				continue
			}
			if p.Regexp.Match([]byte(agent)) {
				return p.Data, true
			}
		}
	}
	return nil, false
}





*/
