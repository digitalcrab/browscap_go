package browscap_go

import (
	"fmt"
	"reflect"
	"strings"
)

type Browser struct {
	built  bool   // has complete data from parents
	parent string //name of the parent

	Browser         string
	BrowserVersion  string
	BrowserMajorVer string
	BrowserMinorVer string
	// Browser, Application, Bot/Crawler, Useragent Anonymizer, Offline Browser,
	// Multimedia Player, Library, Feed Reader, Email Client or unknown
	BrowserType string

	Platform        string
	PlatformShort   string
	PlatformVersion string

	// Mobile Phone, Mobile Device, Tablet, Desktop, TV Device, Console,
	// FonePad, Ebook Reader, Car Entertainment System or unknown
	DeviceType  string
	DeviceName  string
	DeviceBrand string
}

func (browser *Browser) build(browsers map[string]*Browser) {
	if browser.built {
		return
	}

	n := reflect.ValueOf(*browser).NumField()

	current := reflect.ValueOf(browser)
	parent := browser.parent
	for parent != "" {
		b, ok := browsers[parent]
		if !ok {
			break
		}

		parentObj := reflect.ValueOf(b)
		for i := 0; i < n; i++ {
			cField := current.Elem().Field(i)
			if cField.String() != "" {
				continue
			}

			pField := parentObj.Elem().Field(i)
			if pField.String() == "" {
				continue
			}

			cField.SetString(pField.String())
		}

		parent = b.parent
	}

	browser.built = true
}

func (browser *Browser) setValue(key, item string) {
	if key == "Parent" {
		browser.parent = item
	} else if key == "Browser" {
		browser.Browser = item
	} else if key == "Version" {
		browser.BrowserVersion = item
	} else if key == "MajorVer" {
		browser.BrowserMajorVer = item
	} else if key == "MinorVer" {
		browser.BrowserMinorVer = item
	} else if key == "Browser_Type" {
		browser.BrowserType = item
	} else if key == "Platform" {
		browser.Platform = item
		browser.PlatformShort = strings.ToLower(item)

		if strings.HasPrefix(browser.PlatformShort, "win") {
			browser.PlatformShort = "win"
		} else if strings.HasPrefix(browser.PlatformShort, "mac") {
			browser.PlatformShort = "mac"
		}
	} else if key == "Platform_Version" {
		browser.PlatformVersion = item
	} else if key == "Device_Type" {
		browser.DeviceType = item
	} else if key == "Device_Code_Name" {
		browser.DeviceName = item
	} else if key == "Device_Brand_Name" {
		browser.DeviceBrand = item
	}
}

func extractBrowser(data map[string]string) *Browser {
	browser := &Browser{}

	if debug {
		fmt.Println("= Browser ==================")
		for k, v := range data {
			fmt.Printf("%s = %s\n", k, v)
		}
		fmt.Println("============================")
	}

	// Browser
	if item, ok := data["Browser"]; ok {
		browser.Browser = item
	}
	if item, ok := data["Version"]; ok {
		browser.BrowserVersion = item
	}
	if item, ok := data["MajorVer"]; ok {
		browser.BrowserMajorVer = item
	}
	if item, ok := data["MinorVer"]; ok {
		browser.BrowserMinorVer = item
	}
	if item, ok := data["Browser_Type"]; ok {
		browser.BrowserType = item
	}

	// Platform
	if item, ok := data["Platform"]; ok {
		browser.Platform = item
		browser.PlatformShort = strings.ToLower(item)

		if strings.HasPrefix(browser.PlatformShort, "win") {
			browser.PlatformShort = "win"
		} else if strings.HasPrefix(browser.PlatformShort, "mac") {
			browser.PlatformShort = "mac"
		}
	}
	if item, ok := data["Platform_Version"]; ok {
		browser.PlatformVersion = item
	}

	// Device
	if item, ok := data["Device_Type"]; ok {
		browser.DeviceType = item
	}
	if item, ok := data["Device_Code_Name"]; ok {
		browser.DeviceName = item
	}
	if item, ok := data["Device_Brand_Name"]; ok {
		browser.DeviceBrand = item
	}

	return browser
}

func (self *Browser) IsCrawler() bool {
	return self.BrowserType == "Bot/Crawler"
}

func (self *Browser) IsMobile() bool {
	return self.DeviceType == "Mobile Phone" || self.DeviceType == "Mobile Device"
}

func (self *Browser) IsTablet() bool {
	return self.DeviceType == "Tablet" || self.DeviceType == "FonePad" || self.DeviceType == "Ebook Reader"
}

func (self *Browser) IsDesktop() bool {
	return self.DeviceType == "Desktop"
}

func (self *Browser) IsConsole() bool {
	return self.DeviceType == "Console"
}

func (self *Browser) IsTv() bool {
	return self.DeviceType == "TV Device"
}

func (self *Browser) IsAndroid() bool {
	return self.Platform == "Android"
}

func (self *Browser) IsIPhone() bool {
	return self.Platform == "iOS" && self.DeviceName == "iPhone"
}

func (self *Browser) IsIPad() bool {
	return self.Platform == "iOS" && self.DeviceName == "iPad"
}

func (self *Browser) IsWinPhone() bool {
	return strings.Index(self.Platform, "WinPhone") != -1 || self.Platform == "WinMobile"
}
