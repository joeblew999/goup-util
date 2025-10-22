package packaging

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ArchiveFormat represents the archive format to use
type ArchiveFormat string

const (
	TarGz ArchiveFormat = "tar.gz"
	Zip   ArchiveFormat = "zip"
)

// CreateArchive creates an archive of the specified source
func CreateArchive(sourcePath, outputPath string, format ArchiveFormat) error {
	switch format {
	case TarGz:
		return createTarGz(sourcePath, outputPath)
	case Zip:
		return createZip(sourcePath, outputPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// createTarGz creates a tar.gz archive
func createTarGz(sourcePath, outputPath string) error {
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Get base name for relative paths
	baseDir := filepath.Dir(sourcePath)

	// Walk the source directory
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Set relative name
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If not a regular file, skip
		if !info.Mode().IsRegular() {
			return nil
		}

		// Copy file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})
}

// createZip creates a zip archive
func createZip(sourcePath, outputPath string) error {
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Get base name for relative paths
	baseDir := filepath.Dir(sourcePath)

	// Walk the source
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		// Create zip file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		// Create writer for this file
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open source file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy content
		_, err = io.Copy(writer, file)
		return err
	})
}

// CopyFile copies a single file (helper for simple operations)
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
