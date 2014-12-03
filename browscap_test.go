package browscap_go

import (
	"testing"
)

const (
	TEST_INI_FILE     = "./test-data/full_php_browscap.ini"
	TEST_USER_AGENT   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
	TEST_IPHONE_AGENT = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"
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

func TestGetBrowserIPhone(t *testing.T) {
	if browser, ok := GetBrowser(TEST_IPHONE_AGENT); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceName != "iPhone" {
		t.Errorf("Expected iPhone but got %q", browser.DeviceName)
	} else if browser.IsMobile() != true {
		t.Errorf("Expected true but got %t", browser.IsMobile())
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

func TestLastVersion(t *testing.T) {
	version, err := LastVersion()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if version == "" {
		t.Fatalf("Version not found")
	}
	t.Logf("Last version is %q, current version: %q", version, InitializedVersion())
}

func TestDownload(t *testing.T) {
	version, err := LastVersion()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if version != InitializedVersion() {
		t.Logf("Start downloading version %q", version)
		tmpFile := "/tmp/browscap_go_TestDownload.ini"
		if err = DownloadFile(tmpFile); err != nil {
			t.Fatalf("%v", err)
		}

		t.Logf("Initializing with new version")
		InitBrowsCap(tmpFile, true)

		if version != InitializedVersion() {
			t.Fatalf("New file is wrong")
		}
	}
}
