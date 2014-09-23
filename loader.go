package browscap_go

import (
	"bytes"
	"bufio"
	"os"
	"io"
)

var (
	// Ini
	sEmpty		= []byte{}		// empty signal
	nComment	= []byte{'#'}	// number signal
	sComment	= []byte{';'}	// semicolon signal
	sStart		= []byte{'['}	// section start signal
	sEnd		= []byte{']'}	// section end signal
	sEqual		= []byte{'='}	// equal signal
	sQuote1		= []byte{'"'}	// quote " signal
	sQuote2		= []byte{'\''}	// quote ' signal

	// To reduce memory usage we will keep only next keys
	keepKeys	= [][]byte{
		// Required
		[]byte{'P','a','r','e','n','t'},

		// Used in Browser
		[]byte{'B','r','o','w','s','e','r'},
		[]byte{'V','e','r','s','i','o','n'},
		[]byte{'M','a','j','o','r','V','e','r'},
		[]byte{'M','i','n','o','r','V','e','r'},
		[]byte{'B','r','o','w','s','e','r','_','T','y','p','e'},
		[]byte{'P','l','a','t','f','o','r','m'},
		[]byte{'P','l','a','t','f','o','r','m','_','V','e','r','s','i','o','n'},
		[]byte{'D','e','v','i','c','e','_','T','y','p','e'},
		[]byte{'D','e','v','i','c','e','_','C','o','d','e','_','N','a','m','e'},
		[]byte{'D','e','v','i','c','e','_','B','r','a','n','d','_','N','a','m','e'},
	}
)

func loadFromIniFile(path string) (*dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dict := newDictionary()

	buf := bufio.NewReader(file)
	sectionName := ""
	sectionPrefix := ""

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
			sectionPrefix = getPrefix(sectionName)
			continue
		}

		// Create section
		if _, ok := dict.mapped[sectionName]; !ok {
			// Save mapped
			dict.mapped[sectionName] = make(section)

			// Create prefix for section
			if _, exists := dict.expressions[sectionPrefix]; !exists {
				dict.expressions[sectionPrefix] = []*expression{}
			}

			// Build expression
			var ee *expression
			if bytes.IndexAny([]byte(sectionName), "*?") != -1 {
				ee = newRegexpExpression(sectionName)
			} else {
				ee = newCompareExpression(sectionName)
			}
			dict.expressions[sectionPrefix] = append(dict.expressions[sectionPrefix], ee)
		}

		// Key => Value
		kv := bytes.SplitN(line, sEqual, 2)

		// Parse Key
		key := bytes.TrimSpace(kv[0])
		if !inList(key, keepKeys) {
			continue
		}

		// Parse Value
		val := bytes.TrimSpace(kv[1])
		if bytes.HasPrefix(val, sQuote1) {
			val = bytes.Trim(val, `"`)
		}
		if bytes.HasPrefix(val, sQuote2) {
			val = bytes.Trim(val, `'`)
		}

		dict.mapped[sectionName][string(key)] = string(val)
	}

	return dict, nil
}
