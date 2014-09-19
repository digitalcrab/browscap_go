package main

import (
	"os"
	"io"
	"bufio"
	"regexp"
	"fmt"
	"net/http"
	"time"
	"crypto/md5"
)

var c uint64
var e uint64
var t time.Duration
var h map[string]bool

func main() {
	h = make(map[string]bool)
	file, err := os.Open("/Users/fromYukki/Downloads/alice")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	client := &http.Client{}
	buf := bufio.NewReader(file)
	re := regexp.MustCompile(`(?is).*\[.*\].*".*".*"(.*)"`)
	// $http_x_forwarded_for $remote_addr $remote_user [$time_local] $http_x_forwarded_proto "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		res := re.FindSubmatch(line)
		if len(res) < 2 {
			continue
		}
		agent := string(res[1])
		if agent == "" || agent == "-" {
			continue
		}
		/*md5 := Md5(agent)
		if _, ok := h[md5]; !ok {
			h[md5] = true
		}*/
		fmt.Printf("%s\n", agent)

		req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", agent)

		start := time.Now()

		response, err := client.Do(req)
		if err != nil {
			e++
			saveTime(start)
			continue
		}
		saveTime(start)
		fmt.Printf("%d\n", response.StatusCode)
		if response.StatusCode != 200 {
			e++
		}
		response.Body.Close()
		c++

		if c == 2000 {
			duration := t/time.Duration(c)
			var durationUnits string
			switch {
			case duration > 2000000:
				durationUnits = "ms"
				duration /= 1000000
			case duration > 1000:
				durationUnits = "μs"
				duration /= 1000
			default:
				durationUnits = "ns"
			}
			fmt.Printf("Count %d, Errors: %d\n", c, e)
			fmt.Printf("Arg: %d %s\n", duration, durationUnits)
			os.Exit(0)
		}
	}

	fmt.Printf("Count %d, Hash: %d\n", c, len(h))
	//fmt.Printf("Count %d, Errors: %d\n", c, e)
	//fmt.Printf("Arg: %f", (float64(t) / float64(c) / float64(1000000)))
}

func saveTime(s time.Time) {
	t += time.Since(s)
}

func Md5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
