// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"bufio"
	"bytes"
	"os"
)

var (
	// Ini
	nComment = []byte{'#'} // number signal
	sComment = []byte{';'} // semicolon signal
	sStart   = []byte{'['} // section start signal
	sEnd     = []byte{']'} // section end signal

	versionSection = "GJK_Browscap_Version"
	versionKey     = "Version"
)

func loadFromIniFile(path string) (*dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dict := newDictionary()

	sectionName := ""

	lineNum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
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
		kvSplit := bytes.IndexByte(line, '=')

		// Parse Key
		keyb := bytes.TrimSpace(line[0:kvSplit])

		// Parse Value
		valb := bytes.TrimSpace(line[kvSplit+1:])
		if len(valb) == 0 {
			continue
		}
		if valb[0] == '"' || valb[0] == '\'' {
			valb = valb[1 : len(valb)-1]
		}

		key := string(keyb)
		val := string(valb)

		if sectionName == versionSection {
			if key == versionKey {
				version = val
			}
			continue
		}

		if _, ok := dict.browsers[sectionName]; !ok {
			dict.tree.Add(sectionName, lineNum)

			browser := &Browser{
				Section: sectionName,
			}
			browser.setValue(key, val)
			dict.browsers[sectionName] = browser
		} else {
			dict.browsers[sectionName].setValue(key, val)
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dict, nil
}
