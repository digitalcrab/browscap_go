// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"bufio"
	"io/ioutil"
	"strings"
	"testing"
	"os"
)

const (
	TEST_INI_FILE     = "./test-data/full_php_browscap.ini"
	TEST_USER_AGENT   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
	TEST_IPHONE_AGENT = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"
)

func initFromTestIniFile(tb testing.TB) {
	if err := InitBrowsCap(TEST_INI_FILE, false); err != nil {
		tb.Fatalf("%v", err)
	}
}

func TestInitBrowsCap(t *testing.T) {
	if err := InitBrowsCap(TEST_INI_FILE, true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestInitBrowsCapFromReader(t *testing.T) {
	file, err := os.Open(TEST_INI_FILE)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	if err := InitBrowsCapFromReader(buf, true); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetBrowser(t *testing.T) {
	initFromTestIniFile(t)
	if browser, ok := GetBrowser(TEST_USER_AGENT); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "Chrome" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	} else if browser.Platform != "MacOSX" {
		t.Errorf("Expected MacOSX but got %q", browser.Platform)
	} else if browser.BrowserVersion != "37.0" {
		t.Errorf("Expected 37.0 but got %q", browser.BrowserVersion)
	} else if browser.RenderingEngineName != "Blink" {
		t.Errorf("Expected Blink but got %q", browser.RenderingEngineName)
	} else if browser.Crawler != "false" {
		t.Errorf("Expected false but got %q", browser.Crawler)
	}
}

func TestGetBrowserIPhone(t *testing.T) {
	initFromTestIniFile(t)
	if browser, ok := GetBrowser(TEST_IPHONE_AGENT); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceName != "iPhone" {
		t.Errorf("Expected iPhone but got %q", browser.DeviceName)
	} else if browser.Platform != "iOS" {
		t.Errorf("Expected iOS but got %q", browser.Platform)
	} else if browser.PlatformVersion != "4.3" {
		t.Errorf("Expected 4.3 but got %q", browser.PlatformVersion)
	} else if browser.IsMobile() != true {
		t.Errorf("Expected true but got %t", browser.IsMobile())
	}
}

func TestGetBrowserYandex(t *testing.T) {
	initFromTestIniFile(t)
	if browser, ok := GetBrowser("Yandex Browser 1.1"); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "Yandex Browser" {
		t.Errorf("Expected Yandex Browser but got %q", browser.Browser)
	} else if browser.IsCrawler() != false {
		t.Errorf("Expected false but got %t", browser.IsCrawler())
	}
}
func TestGetBrowser360Spider(t *testing.T) {
	initFromTestIniFile(t)
	if browser, ok := GetBrowser("360Spider"); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "360Spider" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	} else if browser.IsCrawler() != true {
		t.Errorf("Expected true but got %t", browser.IsCrawler())
	}
}

func TestGetBrowserIssues(t *testing.T) {
	initFromTestIniFile(t)
	// https://github.com/digitalcrab/browscap_go/issues/4
	ua := "Mozilla/5.0 (iPad; CPU OS 5_0_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A405 Safari/7534.48.3"
	if browser, ok := GetBrowser(ua); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceType != "Tablet" {
		t.Errorf("Expected tablet %q", browser.DeviceType)
	}
}

func TestLastVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

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
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

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

func BenchmarkInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitBrowsCap(TEST_INI_FILE, true)
	}
}

func BenchmarkGetBrowser(b *testing.B) {
	initFromTestIniFile(b)
	
	data, err := ioutil.ReadFile("test-data/user_agents_sample.txt")
	if err != nil {
		b.Error(err)
	}

	uas := strings.Split(strings.TrimSpace(string(data)), "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(uas)

		_, ok := GetBrowser(uas[idx])
		if !ok {
			b.Errorf("User agent not recognized: %s", uas[idx])
		}
	}
}
