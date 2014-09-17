package browscap_go

import (
	"testing"
)

func initBrowsCapXml() error {
	return InitBrowsCapXml("./test-data/browscap.xml", false)
}

func TestInitialise(t *testing.T) {
	if err := initBrowsCapXml(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetBrowser(t *testing.T) {
	if err := initBrowsCapXml(); err != nil {
		t.Errorf("%v", err)
	}

	agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"

	if res, ok := GetBrowser(agent); !ok {
		t.Errorf("Browser not found")
	} else if pl, ok := res["Platform"]; !ok || pl != "MacOSX" {
		t.Errorf("Platform not found or %q != %q", "MacOSX", pl)
	} else if br, ok := res["Browser"]; !ok || br != "Chrome" {
		t.Errorf("Browser not found or %q != %q", "Chrome", br)
	}
}
