package models

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

type PowerPoint struct {
	Path         string
	Files        []*zip.File
	Slides       map[string][]byte
	NotesSlides  map[string]string
	Themes       map[string]string
	Images       map[string]string
	Presentation string
}

// Create an identical copy of the original powerpoint file
func (ppt *PowerPoint) Clone(newPath string) error {
	oldFile, err := zip.OpenReader(ppt.Path)
	if err != nil {
		return err
	}

	newFile, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	zw := zip.NewWriter(newFile)
	defer zw.Close()

	for _, file := range oldFile.File {
		if err := zw.Copy(file); err != nil {
			return err
		}
	}

	return nil
}

func ReadPowerPoint(path string) (*PowerPoint, error) {
	var p PowerPoint
	p.Slides = make(map[string][]byte)
	p.Path = path
	f, err := zip.OpenReader(path)

	if err != nil {
		return nil, errors.New("Error opening file" + err.Error())
	}

	defer f.Close()

	p.Files = f.File

	for _, file := range p.Files {
		if strings.Contains(file.Name, "ppt/slides/slide") {
			slideOpen, _ := file.Open()
			defer slideOpen.Close()
			bytes := readCloserToByte(slideOpen)
			p.Slides[file.Name] = bytes
		}
	}

	return &p, nil
}

func readCloserToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
