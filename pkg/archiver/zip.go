package archiver

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cafi-dev/iac-gen/pkg/logging"
	"go.uber.org/zap"
)

type TGZ struct{}

func NewTGZ() Archiver {
	return TGZ{}
}

func (tgz TGZ) Extension() string {
	return ".tgz"
}

func (tgz TGZ) compress(src string, writer io.Writer) error {
	// tar > gzip > filewriter
	zr := zip.NewWriter(writer)

	// walk through every file in the folder
	if err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		logging.GetLogger().Info("crawling", zap.String("file", file))
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

func (tgz TGZ) Compress(filePath, pathToCompress string) error {
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
	if err := tgz.compress(pathToCompress, f); err != nil {
		return fmt.Errorf("failed to tgz compress folder %q: %w", pathToCompress, err)
	}
	return nil
}

// Sanitize archive file pathing from "G305: Zip Slip vulnerability"
// fixing code vulnerability https://snyk.io/research/zip-slip-vulnerability#go
func (tgz TGZ) sanitizeExtractPath(filePath string, destination string) (string, error) {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s: illegal file path", filePath)
	}
	return destpath, nil
}

func (tgz TGZ) decompress(src io.Reader, dst string) error {
	// ungzip
	zr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	// untar
	tr := tar.NewReader(zr)

	// uncompress each element
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		// add dst + re-format slashes according to system
		target, err := tgz.sanitizeExtractPath(header.Name, dst)
		if err != nil {
			continue
		}

		// check the type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			for {
				// https://github.com/securego/gosec/pull/433
				_, err := io.CopyN(f, tr, 1024)
				if err != nil {
					if err == io.EOF {
						break
					}
					f.Close()
					return err
				}
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
	return nil
}

func (tgz TGZ) Decompress(filePath, pathToDecompress string) error {
	// create path to decompress if not exist
	if _, err := os.Stat(pathToDecompress); err != nil {
		if err := os.MkdirAll(pathToDecompress, 0755); err != nil {
			return err
		}
	}
	// open archive file for reading
	reader, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file from path %q: %w", filePath, err)
	}
	defer reader.Close()
	return tgz.decompress(reader, pathToDecompress)
}
