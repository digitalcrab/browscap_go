// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"unicode"
)

const (
	DownloadUrl     = "http://browscap.org/stream?q=PHP_BrowsCapINI"
	CheckVersionUrl = "http://browscap.org/version-number"
)

type BrowsCap struct {
	dict    *dictionary
	version string
}

func (browscap *BrowsCap) InitializedVersion() string {
	return browscap.version
}

func LastVersion() (string, error) {
	response, err := http.Get(CheckVersionUrl)
	if err != nil {
		return "", fmt.Errorf("browscap: error sending request, %v", err)
	}
	defer response.Body.Close()

	// Get body of response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("browscap: error reading the response data of request, %v", err)
	}

	// Check 200 status
	if response.StatusCode != 200 {
		return "", fmt.Errorf("browscap: error unexpected status code %d", response.StatusCode)
	}

	return string(body), nil
}

func DownloadFile(saveAs string) error {
	response, err := http.Get(DownloadUrl)
	if err != nil {
		return fmt.Errorf("browscap: error sending request, %v", err)
	}
	defer response.Body.Close()

	// Get body of response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("browscap: error reading the response data of request, %v", err)
	}

	// Check 200 status
	if response.StatusCode != 200 {
		return fmt.Errorf("browscap: error unexpected status code %d", response.StatusCode)
	}

	if err = ioutil.WriteFile(saveAs, body, os.ModePerm); err != nil {
		return fmt.Errorf("browscap: error saving file, %v", err)
	}

	return nil
}

func (browscap *BrowsCap) GetBrowser(userAgent string) (browser *Browser, ok bool) {
	agent := mapToBytes(unicode.ToLower, userAgent)
	defer bytesPool.Put(agent)

	name := browscap.dict.tree.Find(agent)
	if name == "" {
		return
	}

	browser = browscap.dict.getBrowser(name)
	if browser != nil {
		ok = true
	}

	return
}
