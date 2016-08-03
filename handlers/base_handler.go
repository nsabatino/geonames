package handlers

import (
	"archive/zip"
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/remizovm/geonames/helpers"
	"github.com/remizovm/geonames/types"
)

type BaseHandler struct {
	URL          string
	FileName     string
	IsZipped     bool
	Type         interface{}
	ParseFunc    func([]string) bool
	HeaderLength uint
	Result       interface{}
}

func (bh *BaseHandler) Process() error {
	// download
	data, err := helpers.HTTPGetNew(bh.URL)
	if err != nil {
		return err
	}
	// make temp file
	tempFile := helpers.GetTempPath(bh.URL)
	// write data to temp file
	f, err := helpers.WriteToFile(tempFile, data)
	if err != nil {
		return err
	}
	// defer cleanup
	defer data.Close()
	defer f.Close()
	defer os.Remove(tempFile)
	// scanner to parse from
	var s *bufio.Scanner
	// handle zipped file
	if bh.IsZipped {
		r, err := zip.OpenReader(f.Name())
		if err != nil {
			return err
		}
		defer r.Close()
		// get txt name
		txtName := strings.Replace(bh.FileName, "zip", "txt", -1)
		for i := range r.File {
			if r.File[i].Name == txtName {
				rc, e := r.File[i].Open()
				if e != nil {
					return e
				}
				defer rc.Close()
				// found the file in zip archive - get the scanner
				s = bufio.NewScanner(rc)
				break
			}
		}
	} else {
		// not zipped
		s = bufio.NewScanner(f)
	}

	if s == nil {
		return errors.New("unknown error")
	}
	// determine result type
	switch bh.Type {
	case types.Country{}:
		bh.Result = make(map[int64]*types.Country)
	}
	// ??? parsing
	helpers.StringParse(s, bh.HeaderLength, bh.ParseFunc)

	return nil
}
