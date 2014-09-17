package browscap_go

import (
	"fmt"
	"github.com/fromYukki/browscap_go/ini"
)

var (
	initialized bool
)

func InitBrowsCap(path string, force bool) error {
	if initialized && !force {
		return nil
	}

	var (
		dic	ini.Dictionary
		err	error
	)

	// Load ini file
	if dic, err = ini.LoadFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	// Process ini file
	if err = process(dic); err != nil {
		return fmt.Errorf("browscap: An error occurred while processing data, %v ", err)
	}

	initialized = true
	return nil
}

func process(dic ini.Dictionary) error {
	return nil
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

	regexReplace = map[string]string{
		"(": "\\(",
		")": "\\)",
		"[": "\\[",
		"]": "\\]",
		"{": "\\{",
		"}": "\\}",
		"<": "\\<",
		">": "\\>",
		"$": "\\$",
		"^": "\\^",
		"+": "\\+",
		"!": "\\!",
		"=": "\\=",
		"|": "\\|",
		":": "\\:",
		"-": "\\-",
		//
		"*": ".*",
		"?": ".",
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

func newPattern(pattern string) *Pattern {
	return &Pattern{
		Pattern:	pattern,
		Data:		make(map[string]string),
	}
}

func rebase(patternsMap PatternsMap) (Patterns, error) {
	res := Patterns{}

	for name, pattern := range patternsMap {
		// Collect data items
		resData := findData(patternsMap, name)

		// Rebuild data without parents
		pattern.Data = make(map[string]string)
		if len(resData) > 0 {
			for k, v := range resData {
				if k == "Parent" {
					continue
				}
				pattern.Data[k] = v
			}
		}

		// Build regexp
		exp, err := regexp.Compile("(?i)^" + escapePattern(pattern.Pattern) + "$")
		if err != nil {
			return nil, err
		}
		pattern.Regexp = exp

		res = append(res, pattern)
	}

	return res, nil
}

func findData(patternsMap PatternsMap, name string) (map[string]string) {
	res := make(map[string]string)

	if item, found := patternsMap[name]; found {
		// Parent's data
		if parentName, hasParent := item.Data["Parent"]; hasParent {
			parentData := findData(patternsMap, parentName)
			if len(parentData) > 0 {
				for k, v := range parentData {
					if k == "Parent" {
						continue
					}
					res[k] = v
				}
			}
		}
		// It's item data
		if len(item.Data) > 0 {
			for k, v := range item.Data {
				if k == "Parent" {
					continue
				}
				res[k] = v
			}
		}
	}

	return res
}

func escapePattern(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, ".", "\\.", -1)

	for k, v := range regexReplace {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}*/
