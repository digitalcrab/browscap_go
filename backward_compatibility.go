package browscap_go

import "fmt"

var backwardCompatibleBrowsCap *BrowsCap

// Deprecated: Use NewBrowsCapFromFile
func InitBrowsCap(path string, force bool) error {
	if backwardCompatibleBrowsCap != nil && !force {
		return nil
	}
	var err error

	// Load ini file
	if backwardCompatibleBrowsCap, err = NewBrowsCapFromFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	return nil
}

// Deprecated: Use BrowsCap.InitializedVersion
func InitializedVersion() string {
	return backwardCompatibleBrowsCap.InitializedVersion()
}

// Deprecated: Use BrowsCap.GetBrowser
func GetBrowser(userAgent string) (browser *Browser, ok bool) {
	if backwardCompatibleBrowsCap == nil {
		return
	}

	return backwardCompatibleBrowsCap.GetBrowser(userAgent)
}
