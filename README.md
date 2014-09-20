# Browser Capabilities GoLang Project

PHP has `get_browser()` function which tells what the user's browser is capable of.
You can check original documentation [here](http://php.net/get_browser). 
This is GoLang analog of `get_browser()` function.

[![Build Status](https://secure.travis-ci.org/fromYukki/browscap_go.png?branch=master)](http://travis-ci.org/fromYukki/browscap_go)

## Introduction

The [browscap.ini](http://browscap.org/) file is a database which provides a lot of details about 
browsers and their capabilities, such as name, versions, Javascript support and so on.

## Quick start

First of all you need initialize library with [browscap.ini](http://browscap.org/) file. 
And then you can get Browser information as `Browser` structure.

```
import (
	"fmt"
	bgo "github.com/fromYukki/browscap_go"
)

func main() {
	if err := bgo.InitBrowsCap("browscap.ini", false); err != nil {
		panic(err)
	}
	
	browser, ok := bgo.GetBrowser("")
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
