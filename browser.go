package browscap_go

type Browser struct {
	Browser			string
	BrowserVersion	string
	BrowserType		string
	Platform		string
	PlatformVersion	string
	DeviceType		string
	DeviceName		string
	DeviceBrand		string
	IsCrawler		bool
	IsMobileDevice	bool
	IsTablet		bool
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
	if item, ok := data["Crawler"]; ok {
		browser.IsCrawler = item == "true"
	}
	if item, ok := data["isMobileDevice"]; ok {
		browser.IsMobileDevice = item == "true"
	}
	if item, ok := data["isTablet"]; ok {
		browser.IsTablet = item == "true"
	}

	return browser
}

