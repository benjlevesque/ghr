package util

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/h2non/filetype"
)

// ExtractTarGzBinary extracts only the binary content of a tar gz tothe target directory
func ExtractTarGzBinary(gzipStream io.Reader, target string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeReg:
			data, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return err
			}
			kind, err := filetype.MatchReader(bytes.NewReader(data))
			if err != nil {
				return err
			}

			if kind.MIME.Type == "application" {
				outFile, err := os.OpenFile(target+"/"+header.Name, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				defer outFile.Close()
				err = ioutil.WriteFile(target+"/"+header.Name, data, os.FileMode(header.Mode))
				// }
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return errors.New("No application found")
}
