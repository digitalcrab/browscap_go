package browscap_go

import (
	"fmt"
	"bytes"
)

var (
	dict		*dictionary
	initialized bool
	debug		bool
)

func Debug(val bool) {
	debug = val
}

func InitBrowsCap(path string, force bool) error {
	if initialized && !force {
		return nil
	}
	var err	error

	// Load ini file
	if dict, err = loadFromIniFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	initialized = true
	return nil
}

func GetBrowser(userAgent string) (browser *Browser, ok bool) {
	if !initialized {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			browser = nil
			ok = false
		}
	}()

	agent := bytes.ToLower([]byte(userAgent))
	prefix := getPrefix(userAgent)

	// Main search
	if browser, ok = getBrowser(prefix, agent); ok {
		return
	}

	// Fallback
	if prefix != "*" {
		browser, ok = getBrowser("*", agent)
	}

	return
}

func getBrowser(prefix string, agent []byte) (browser *Browser, ok bool) {
	if expressions, exists := dict.expressions[prefix]; exists {
		for _, exp := range expressions {
			if exp.Match(agent) {
				data := dict.findData(exp.Name)
				browser = extractBrowser(data)
				ok = true
				return
			}
		}
	}
	return
}
