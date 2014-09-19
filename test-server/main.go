package main

import (
	"fmt"
	"net/http"
	"github.com/fromYukki/browscap_go"
	"runtime"
)

func main() {
	err := browscap_go.InitBrowsCap("../test-data/full_php_browscap.ini", false)
	if err != nil {
		panic(err)
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Printf("%d bytes\n", ms.Alloc)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			if browser, ok := browscap_go.GetBrowser(r.UserAgent()); ok {
				w.WriteHeader(200)
				fmt.Fprintf(w, "%s %s", browser.Browser, browser.Platform)
			} else {
				w.WriteHeader(404)
			}
			fmt.Printf("%d bytes\n", ms.Alloc)
		})
	http.ListenAndServe(":8080", nil)
}
