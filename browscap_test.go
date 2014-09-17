package browscap_go

import (
	"testing"
	"github.com/fromYukki/browscap_go/ini"
)

const (
	TEST_INI_FILE = "./test-data/full_php_browscap.ini"
)

func TestIni(t *testing.T) {
	dict, err := ini.LoadFile(TEST_INI_FILE)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("Size: %d", len(dict))
}
