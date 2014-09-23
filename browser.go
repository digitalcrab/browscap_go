package browscap_go

import (
	"strings"Fix
	"fmt"
)

type Browser struct {
	Browser			string
	BrowserVersion	string
	BrowserMajorVer	string
	BrowserMinorVer	string
	// Browser, Application, Bot/Crawler, Useragent Anonymizer, Offline Browser,
	// Multimedia Player, Library, Feed Reader, Email Client or unknown
	BrowserType		string

	Platform		string
	PlatformShort	string
	PlatformVersion	string

	// Mobile Phone, Mobile Device, Tablet, Desktop, TV Device, Console,
	// FonePad, Ebook Reader, Car Entertainment System or unknown
	DeviceType		string
	DeviceName		string
	DeviceBrand		string
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
