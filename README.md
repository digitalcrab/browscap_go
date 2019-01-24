Golang Client for Browser Capabilities Project
==============================================

[Browscap](http://browscap.org/) provides regularly updated rules for parsing user agents.

For more info check the [officall PHP client](https://github.com/browscap/browscap-php).

Example:
--------

```go
import (
	"fmt"
	bgo "github.com/digitalcrab/browscap_go"
)

func main() {
	if err := bgo.InitBrowsCap("browscap.ini", false); err != nil {
		panic(err)
	}
	
	browser, ok := bgo.GetBrowser("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36")
	if !ok || browser == nil {
    	panic("Browser not found")
	} else {
    	fmt.Printf("Browser = %s [%s] v%s\n", browser.Browser, browser.BrowserType, browser.BrowserVersion)
    	fmt.Printf("Platform = %s v%s\n", browser.Platform, browser.PlatformVersion)
    	fmt.Printf("Device = %s [%s] %s\n", browser.DeviceName, browser.DeviceType, browser.DeviceBrand)
    	fmt.Printf("IsCrawler = %t\n", browser.IsCrawler())
    	fmt.Printf("IsMobile = %t\n", browser.IsMobile())
	}
}
```

Browser object fields
---------------------

```go
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

Crawler string

Cookies    string
JavaScript string

RenderingEngineName    string
RenderingEngineVersion string
```

Bechmark
--------

```
BenchmarkInit-4         	       1	2256895168 ns/op	346136904 B/op	 5700912 allocs/op
BenchmarkGetBrowser-4   	   10000	    140975 ns/op	      37 B/op	       1 allocs/op
```

## License

```
The MIT License (MIT)

Copyright (c) 2015 Maksim Naumov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
