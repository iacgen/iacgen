package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZIP struct{}

func NewZIP() Archiver {
	return ZIP{}
}

func (z ZIP) Extension() string {
	return ".zip"
}

func (z ZIP) compress(src string, writer io.Writer) error {
	zr := zip.NewWriter(writer)

	// walk through every file in the folder
	if err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		// must provide real name
		filename := strings.TrimPrefix(strings.TrimPrefix(file, src), "/")
		fw, err := zr.Create(filename)
		if err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		// copy over contents
		for {
			// https://github.com/securego/gosec/pull/433
			_, err := io.CopyN(fw, f, 1024)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// produce zip
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}

func (z ZIP) Compress(filePath, pathToCompress string) error {
	// create path to compress if not exist
	fileDir := filepath.Dir(filePath)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		// dir does not exist
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for file %q: %w", filePath, err)
		}
	}
	// create archive file
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("failed to open file to write %q: %w", filePath, err)
	}
	defer f.Close()
	if err := z.compress(pathToCompress, f); err != nil {
		return fmt.Errorf("failed to zip folder %q: %w", pathToCompress, err)
	}
	return nil
}
