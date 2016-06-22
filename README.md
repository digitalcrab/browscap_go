# Browser Capabilities GoLang Project

PHP has `get_browser()` function which tells what the user's browser is capable of.
You can check original documentation [here](http://php.net/get_browser). 
This is GoLang analog of `get_browser()` function.

[![Build Status](https://secure.travis-ci.org/digitalcrab/browscap_go.png?branch=master)](http://travis-ci.org/digitalcrab/browscap_go)

## Introduction

The [browscap.ini](http://browscap.org/) file is a database which provides a lot of details about 
browsers and their capabilities, such as name, versions, Javascript support and so on.

## Quick start

First of all you need initialize library with [browscap.ini](http://browscap.org/) file. 
And then you can get Browser information as `Browser` structure.

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
