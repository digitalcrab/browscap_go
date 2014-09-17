package browscap_go

/*
import (
	"encoding/xml"
	"io/ioutil"
	"fmt"
)

type xmlBrowserCaps struct {
	Version		xmlBrowserVersion
	Patterns	[]xmlBrowserPattern	`xml:"browsercapitems>browscapitem"`
}

type xmlBrowserVersion struct {
	XMLName		xml.Name   			`xml:"gjk_browscap_version"`
	Items		[]xmlBrowserItem	`xml:"item"`
}

type xmlBrowserItem struct {
	Name	string	`xml:"name,attr"`
	Value	string	`xml:"value,attr"`
}

type xmlBrowserPattern struct {
	XMLName	xml.Name   			`xml:"browscapitem"`
	Pattern	string				`xml:"name,attr"`
	Items	[]xmlBrowserItem	`xml:"item"`
}

func parseXml(path string) (*xmlBrowserCaps, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("browscap: unable open xml file, %v", err)
	}

	var caps xmlBrowserCaps
	if err = xml.Unmarshal(buf, &caps); err != nil {
		return nil, fmt.Errorf("browscap: error on paring xml file, %v", err)
	}

	return &caps, nil
}

func buildFromXml(caps *xmlBrowserCaps) PatternsMap {
	patterns := make(PatternsMap)

	for _, xmlPattern := range caps.Patterns {
		pattern := newPattern(xmlPattern.Pattern)
		for _, item := range xmlPattern.Items {
			pattern.Data[item.Name] = item.Value
		}
		patterns[xmlPattern.Pattern] = pattern
	}

	return patterns
}
*/
