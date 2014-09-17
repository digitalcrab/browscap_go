package browscap_go

import (
	"bytes"
	"bufio"
	"os"
	"io"
)

var (
	sEmpty		= []byte{}		// empty signal
	nComment	= []byte{'#'}	// number signal
	sComment	= []byte{';'}	// semicolon signal
	sStart		= []byte{'['}	// section start signal
	sEnd		= []byte{']'}	// section end signal
	sEqual		= []byte{'='}	// equal signal
	sQuote1		= []byte{'"'}	// quote " signal
	sQuote2		= []byte{'\''}	// quote ' signal
)

func loadFromIniFile(path string) (*dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Store parsed section and it's indexes
	sections := make(map[string]int)
	dict := newDictionary()

	buf := bufio.NewReader(file)
	section := ""

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
			section = string(line[1 : len(line)-1])
			continue
		}

		// Create section
		idx, ok := sections[section]
		if !ok {
			// Append to sorted list
			dict.sorted = append(dict.sorted, newSection(section))
			// Save index
			idx = len(dict.sorted) - 1
			// Save parsed
			sections[section] = idx
			// Save mapped
			dict.mapped[section] = dict.sorted[idx]
		}

		// Key => Value
		kv := bytes.SplitN(line, sEqual, 2)

		// Parse Value
		val := bytes.TrimSpace(kv[1])
		if bytes.HasPrefix(val, sQuote1) {
			val = bytes.Trim(val, `"`)
		}
		if bytes.HasPrefix(val, sQuote2) {
			val = bytes.Trim(val, `'`)
		}

		// Parse Key
		key := string(bytes.TrimSpace(kv[0]))
		dict.sorted[idx].Data[key] = string(val)
	}

	return dict, nil
}
