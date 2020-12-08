package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

// StringCounter struct to counting substring in data,
// data reads from io.Reader
type StringCounter struct {
	substring string
	reader    io.Reader
}

// SetSubstring - getter to substr in  StringCounter
func (p *StringCounter) SetSubstring(aSubstring string) {
	p.substring = aSubstring
}

// GetSubstring - getter to substr in  StringCounter
func (p *StringCounter) GetSubstring() string {
	return p.substring
}

// SetReader - setter to reader in  StringCounter
func (p *StringCounter) SetReader(aReader io.Reader) {
	p.reader = aReader
}

// GetReader - getter to reader in  StringCounter
func (p *StringCounter) GetReader() io.Reader {
	return p.reader
}

// SafetyCount func to counting substr in Reader (use buffers)
func (p *StringCounter) SafetyCount() (int, error) {
	var resultCount int = 0

	if p.reader == nil {
		return resultCount, errors.New("reader ref is nil")
	}
	if p.substring == "" {
		return resultCount, errors.New("empty substring")
	}

	sizeOfSubstring := len([]byte(p.substring))
	sliceSize := 1024
	sliceToSearch := []byte(p.substring)

	if len(p.substring) > sliceSize {
		sliceSize = (int)((float32)(len(p.substring)) * 1.1)
	}

	sliceToParse1 := make([]byte, sliceSize, sliceSize)
	sliceToParse2 := make([]byte, sliceSize, sliceSize)

	var err error = nil
	var n1, n2 int

	n1, err = p.reader.Read(sliceToParse1)

	if err != nil {
		if err == io.EOF {
			return 0, nil
		}
		return 0, err
	}

	needPostWork := false
	for {

		if n1 > 0 {
			resultCount += bytes.Count(sliceToParse1, sliceToSearch)
		}
		n2, err = p.reader.Read(sliceToParse2)

		if err != nil && err != io.EOF {
			return 0, err
		}

		if n2 > 0 {
			//joint parse counting
			//example: pattern:go,buf 3; sl1=gog;sl2=og; need catch
			//substring inside joint
			startPosInSlice1 := len(sliceToParse1) - (sizeOfSubstring - 1)
			endPosInSLice2 := 0
			if n2 > sizeOfSubstring-1 {
				endPosInSLice2 = sizeOfSubstring - 1
			} else {
				endPosInSLice2 = n2
			}
			jointSlice := make([]byte, 0, 2*sizeOfSubstring-2)
			jointSlice = append(jointSlice, sliceToParse1[startPosInSlice1:]...)
			jointSlice = append(jointSlice, sliceToParse2[:endPosInSLice2]...)

			resultCount += bytes.Count(jointSlice, sliceToSearch)
		}

		needPostWork = true
		sliceToParse1 = sliceToParse2
		n1 = n2
		sliceToParse2 = make([]byte, sliceSize, sliceSize)
		if err == io.EOF {
			break
		}
	}

	if needPostWork {
		resultCount += bytes.Count(sliceToParse1, sliceToSearch)
	}

	return resultCount, nil
}

// Count func to counting substr in Reader (use ioutill)
func (p *StringCounter) Count() (int, error) {
	var resultCount int = 0

	if p.reader == nil {
		return resultCount, errors.New("reader ref is nil")
	}
	if p.substring == "" {
		return resultCount, errors.New("empty substring")
	}

	s1, err := ioutil.ReadAll(p.reader)

	if err != nil {
		return 0, err
	}

	resultCount = bytes.Count(s1, []byte(p.substring))

	return resultCount, nil
}
