// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

var (
	// Ini
	sEmpty   = []byte{}     // empty signal
	nComment = []byte{'#'}  // number signal
	sComment = []byte{';'}  // semicolon signal
	sStart   = []byte{'['}  // section start signal
	sEnd     = []byte{']'}  // section end signal
	sEqual   = []byte{'='}  // equal signal
	sQuote1  = []byte{'"'}  // quote " signal
	sQuote2  = []byte{'\''} // quote ' signal

	versionSection = "GJK_Browscap_Version"
	versionKey     = "Version"
)

func loadFromIniFile(path string) (*dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	return loadFromReader(buf)
}

func loadFromReader(buf *bufio.Reader) (*dictionary, error) {
	dict := newDictionary()

	sectionName := ""

	lineNum := 0

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		// Empty line
		if bytes.Equal(sEmpty, line) {
			continue
		}

		// Trim
		line = bytes.TrimSpace(line)

		// Empty line
		if bytes.Equal(sEmpty, line) {
			continue
		}

		// Comment line
		if bytes.HasPrefix(line, nComment) || bytes.HasPrefix(line, sComment) {
			continue
		}

		// Section line
		if bytes.HasPrefix(line, sStart) && bytes.HasSuffix(line, sEnd) {
			sectionName = string(line[1 : len(line)-1])
			continue
		}

		// Key => Value
		kv := bytes.SplitN(line, sEqual, 2)

		// Parse Key
		keyb := bytes.TrimSpace(kv[0])

		// Parse Value
		valb := bytes.TrimSpace(kv[1])
		if bytes.HasPrefix(valb, sQuote1) {
			valb = bytes.Trim(valb, `"`)
		}
		if bytes.HasPrefix(valb, sQuote2) {
			valb = bytes.Trim(valb, `'`)
		}

		key := string(keyb)
		val := string(valb)

		if sectionName == versionSection {
			if key == versionKey {
				version = val
			}
			continue
		}

		// Create section
		if _, ok := dict.browsers[sectionName]; !ok {
			dict.tree.Add(sectionName, lineNum)
			dict.browsers[sectionName] = &Browser{}
			lineNum++
		}

		dict.browsers[sectionName].setValue(key, val)
	}

	return dict, nil
}
