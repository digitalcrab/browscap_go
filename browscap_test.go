package browscap_go

import (
	"testing"
)

const (
	TEST_INI_FILE = "./test-data/full_php_browscap.ini"
	TEST_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
)

func TestLoadIni(t *testing.T) {
	dict, err := loadFromIniFile(TEST_INI_FILE)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("Size: %d", len(dict.sorted))
}

func TestInitBrowsCap(t *testing.T) {
	if err := InitBrowsCap(TEST_INI_FILE, false); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetBrowser(t *testing.T) {
	if browser, ok := GetBrowser(TEST_USER_AGENT); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "Chrome" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	} else if browser.Platform != "MacOSX" {
		t.Errorf("Expected MacOSX but got %q", browser.Platform)
	}
}
