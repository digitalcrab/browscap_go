package browscap_go

import (
	"fmt"
	"bytes"
)

var (
	dict		*dictionary
	initialized bool
	//iterations	int
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

	//iterations = 0

	agent1 := []byte(userAgent)
	agent2 := bytes.ToLower(agent1)
	prefix := getPrefix(userAgent)

	// Main search
	if browser, ok = getBrowser(prefix, agent1, agent2); ok {
		return
	}

	// Fallback
	if prefix != "*" {
		browser, ok = getBrowser("*", agent1, agent2)
	}

	return
}

func getBrowser(prefix string, agent1, agent2 []byte) (browser *Browser, ok bool) {
	if expressions, exists := dict.expressions[prefix]; exists {
		//fmt.Printf("Len: %d\n", len(expressions))
		for _, exp := range expressions {
			//iterations++
			if exp.Match(agent1, agent2) {
				data := dict.findData(exp.Name)
				browser = extractBrowser(data)
				ok = true
				return
			}
		}
	}
	return
}
