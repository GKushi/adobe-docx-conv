package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipManager struct{}

func (zm *ZipManager) UnZipFile(path string, distPath string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("error unzipping path: %w", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(distPath, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(distPath)+string(os.PathSeparator)) {
			return fmt.Errorf("error unzipping path: %w", err)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("error unzipping path: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("error unzipping path: %w", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("error unzipping path: %w", err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("error unzipping path: %w", err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("error unzipping path: %w", err)
		}

		defer dstFile.Close()
		defer fileInArchive.Close()
	}
	return nil
}

func (zm *ZipManager) ZipDir(path string) error {
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error zipping path: %w", err)
	}

	zipFileName := path + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("error zipping path: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error zipping path: %w", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error zipping path: %w", err)
		}
		defer file.Close()

		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return fmt.Errorf("error zipping path: %w", err)
		}

		if !fileInfo.IsDir() {
			zipEntry, err := zipWriter.Create(relPath)
			if err != nil {
				return fmt.Errorf("error zipping path: %w", err)
			}

			_, err = io.Copy(zipEntry, file)
			if err != nil {
				return fmt.Errorf("error zipping path: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error zipping path: %w", err)
	}

	return nil
}
