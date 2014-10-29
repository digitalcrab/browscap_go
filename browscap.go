package browscap_go

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	DownloadUrl     = "http://browscap.org/stream?q=PHP_BrowsCapINI"
	CheckVersionUrl = "http://browscap.org/version-number"
)

var (
	dict        *dictionary
	initialized bool
	version     string
	debug       bool
)

func Debug(val bool) {
	debug = val
}

func InitBrowsCap(path string, force bool) error {
	if initialized && !force {
		return nil
	}
	var err error

	// Load ini file
	if dict, err = loadFromIniFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	if verDictionary, exists := dict.mapped["GJK_Browscap_Version"]; exists {
		version = verDictionary["Version"]
	}

	initialized = true
	return nil
}

func InitializedVersion() string {
	return version
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

func GetBrowser(userAgent string) (browser *Browser, ok bool) {
	if !initialized {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			browser = nil
			ok = false
		}
	}()

	agent := bytes.ToLower([]byte(userAgent))
	prefix := getPrefix(userAgent)

	// Main search
	if browser, ok = getBrowser(prefix, agent); ok {
		return
	}

	// Fallback
	if prefix != "*" {
		browser, ok = getBrowser("*", agent)
	}

	return
}

func getBrowser(prefix string, agent []byte) (browser *Browser, ok bool) {
	if expressions, exists := dict.expressions[prefix]; exists {
		for _, exp := range expressions {
			if exp.Match(agent) {
				data := dict.findData(exp.Name)
				browser = extractBrowser(data)
				ok = true
				return
			}
		}
	}
	return
}
