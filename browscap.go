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
