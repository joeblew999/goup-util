package installer

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract unpacks a given archive file to the specified destination directory.
func Extract(source, destination string) error {
	archiveType, err := getArchiveType(source)
	if err != nil {
		return err
	}

	switch archiveType {
	case "zip":
		return unzip(source, destination)
	case "tar.gz":
		return untarGz(source, destination)
	default:
		return fmt.Errorf("unsupported archive format: %s", archiveType)
	}
}

func getArchiveType(source string) (string, error) {
	file, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the first few bytes to determine the archive type
	buffer := make([]byte, 4)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Check for zip magic number
	if bytes.Equal(buffer, []byte{0x50, 0x4b, 0x03, 0x04}) {
		return "zip", nil
	}

	// Check for gzip magic number
	if bytes.Equal(buffer[:2], []byte{0x1f, 0x8b}) {
		return "tar.gz", nil
	}

	// Fallback for xip files
	if strings.HasSuffix(source, ".xip") {
		return "zip", nil
	}

	return "", fmt.Errorf("unable to determine archive type for %s", source)
}

func untarGz(source, destination string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			return fmt.Errorf("unknown type: %c in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

func unzip(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		fpath := filepath.Join(destination, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
