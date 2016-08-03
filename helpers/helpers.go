package helpers

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	// GeonamesURL is the main geonames dump url
	GeonamesURL   = "http://download.geonames.org/export/dump/"
	commentSymbol = byte('#')
	newLineSymbol = byte('\n')
	delimSymbol   = byte('\t')
	// BoolTrue ...
	BoolTrue = "1"
)

// GetTempPath ...
func GetTempPath(name string) string {
	tempDir := os.TempDir()
	return path.Join(tempDir, name)
}

// WriteToFile ...
func WriteToFile(fileName string, data io.ReadCloser) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, data)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getRaw(url, name string) (*bufio.Scanner, error) {
	var err error
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	lowName := strings.ToLower(name)

	isZip := false
	if strings.Contains(lowName, "zip") {
		isZip = true
	}

	tempDir := os.TempDir()

	filePath := path.Join(tempDir, name)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer os.Remove(filePath)

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	if written != resp.ContentLength {
		errMsg := fmt.Sprintf("%s %d %d", file.Name(), written, resp.ContentLength)
		return nil, errors.New(errMsg)
	}

	var result *bufio.Scanner
	if isZip {
		r, e := zip.OpenReader(file.Name())
		if e != nil {
			return nil, e
		}
		defer r.Close()
		txtName := strings.Replace(name, "zip", "txt", -1)
		for i := range r.File {
			if r.File[i].Name == txtName {
				rc, e := r.File[i].Open()
				if e != nil {
					return nil, e
				}
				defer rc.Close()

				result = bufio.NewScanner(rc)
				break
			}
		}
	} else {
		result = bufio.NewScanner(file)
	}

	return result, nil
}

// HTTPGet returns contents of url in a byte slice
func HTTPGet(url string) ([]byte, error) {
	var err error
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// HTTPGetNew ...
func HTTPGetNew(url string) (io.ReadCloser, error) {
	var err error
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Unzip ...
func Unzip(data []byte) ([]*zip.File, error) {
	var err error

	r, err := zip.NewReader(bytes.NewReader(data), (int64)(len(data)))
	if err != nil {
		return nil, err
	}

	return r.File, nil
}

// GetZipData ...
func GetZipData(files []*zip.File, name string) ([]byte, error) {
	var result []byte

	for _, f := range files {
		if f.Name == name {
			src, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer src.Close()

			result, err = ioutil.ReadAll(src)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// StringParse ...
func StringParse(s *bufio.Scanner, headerLength uint, f func([]string) bool) {
	var err error
	var line string
	var rawSplit []string
	for s.Scan() {
		if headerLength != 0 {
			headerLength--
			continue
		}
		line = s.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == commentSymbol {
			continue
		}
		rawSplit = strings.Split(line, "\t")
		if !f(rawSplit) {
			break
		}
	}
	if err = s.Err(); err != nil {
		log.Fatal(err)
	}
}

// Parse ...
func Parse(data []byte, headerLength int, f func([][]byte) bool) {
	rawSplit := bytes.Split(data, []byte{newLineSymbol})
	var rawLineSplit [][]byte
	for i := range rawSplit {
		if headerLength != 0 {
			headerLength--
			continue
		}
		if len(rawSplit[i]) == 0 {
			continue
		}
		if rawSplit[i][0] == commentSymbol {
			continue
		}
		rawLineSplit = bytes.Split(rawSplit[i], []byte{delimSymbol})
		if !f(rawLineSplit) {
			break
		}
	}
}
