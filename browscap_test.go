package browscap_go

import (
	"testing"
)

const (
	TEST_INI_FILE = "./test-data/full_php_browscap.ini"
	TEST_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
)

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

func TestGetBrowserYandex(t *testing.T) {
	if browser, ok := GetBrowser("Yandex Browser 1.1"); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "Yandex Browser" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	}
}

func TestGetBrowser360Spider(t *testing.T) {
	if browser, ok := GetBrowser("360Spider"); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "360Spider" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	}
}
