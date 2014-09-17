package ini

import (
	"bytes"
	"bufio"
	"os"
	"io"
)

const (
	DEFAULT_SESSION_NAME = "default"
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

// Slice to keep order of sections
type Dictionary []*Section

type Section struct {
	Name	string
	Data	map[string]string
}

func LoadFile(path string) (Dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Store parsed section and it's indexes
	sections := make(map[string]int)
	dict := Dictionary{}

	buf := bufio.NewReader(file)
	section := DEFAULT_SESSION_NAME

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
			dict = append(dict, newSection(section))
			idx = len(dict) - 1
			sections[section] = idx
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
		dict[idx].Data[key] = string(val)
	}

	return dict, nil
}

func newSection(name string) *Section {
	return &Section{
		Name:	name,
		Data:	make(map[string]string),
	}
}
