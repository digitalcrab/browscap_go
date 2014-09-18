package browscap_go

type Browser struct {
	Browser			string
	BrowserVersion	string
	// Browser, Application, Bot/Crawler, Useragent Anonymizer, Offline Browser,
	// Multimedia Player, Library, Feed Reader, Email Client or unknown
	BrowserType		string
	Platform		string
	PlatformVersion	string
	// Mobile Phone, Mobile Device, Tablet, Desktop, TV Device, Console,
	// FonePad, Ebook Reader, Car Entertainment System or unknown
	DeviceType		string
	DeviceName		string
	DeviceBrand		string
}

func extractBrowser(data map[string]string) *Browser {
	browser := &Browser{}

	if item, ok := data["Browser"]; ok {
		browser.Browser = item
	}
	if item, ok := data["Version"]; ok {
		browser.BrowserVersion = item
	}
	if item, ok := data["Browser_Type"]; ok {
		browser.BrowserType = item
	}
	if item, ok := data["Platform"]; ok {
		browser.Platform = item
	}
	if item, ok := data["Platform_Version"]; ok {
		browser.PlatformVersion = item
	}
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
