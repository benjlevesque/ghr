package util

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
)

// ExtractTarGzBinary extracts only the binary content of a tar gz tothe target directory
func ExtractTarGzBinary(gzipStream io.Reader, target string) (string, error) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return "", err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch header.Typeflag {
		case tar.TypeReg:
			data, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return "", err
			}
			kind, err := filetype.MatchReader(bytes.NewReader(data))
			if err != nil {
				return "", err
			}

			if kind.MIME.Type == "application" {
				path := target + "/" + filepath.Base(header.Name)
				outFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return "", err
				}
				defer outFile.Close()
				err = ioutil.WriteFile(path, data, os.FileMode(header.Mode))
				// }
				if err != nil {
					return "", err
				}
				return path, nil
			}
		}
	}
	return "", errors.New("No application found")
}
