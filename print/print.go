package print

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/alexbrainman/printer"
)

func AddFileByPath(name, path string) error {

	p, err := printer.Open(name)
	if err != nil {
		return fmt.Errorf("failed open printer '%s': %w", name, err)
	}
	defer p.Close()

	splitted := strings.Split(path, ".")
	extention := splitted[len(splitted)-1]

	switch extention {
	case "jpg":
		err = addImg(p, path)
		if err != nil {
			return fmt.Errorf("failed adding img to '%s': %w", name, err)
		}
		return nil
	default:
		return errors.New("file extention is not supported")
	}
}

func addImg(p *printer.Printer, img string) error {
	err := p.StartPage()
	if err != nil {
		return fmt.Errorf("failed starting page: %w", err)
	}

	// Open the image file
	f, err := os.Open(img)
	if err != nil {
		return fmt.Errorf("failed opening img: %w", err)
	}
	defer f.Close()

	// Decode the image
	data, err := jpeg.Decode(f)
	if err != nil {
		return fmt.Errorf("failed decoding img: %w", err)
	}

	// Encode the image as JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, data, nil)
	if err != nil {
		return fmt.Errorf("failed encoding img: %w", err)
	}

	// Write the data to the printer
	_, err = p.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed writing img to printer: %w", err)
	}

	// End the page
	err = p.EndPage()
	if err != nil {
		return fmt.Errorf("failed ending page: %w", err)
	}

	return nil
}
