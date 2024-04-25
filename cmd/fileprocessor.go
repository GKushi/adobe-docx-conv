package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileProcessor struct {
	DocxConverter *DocxConverter
	ZipManager    *ZipManager
}

func (fp *FileProcessor) Process(pathArg string) error {
	absolutePath, err := filepath.Abs(pathArg)
	if err != nil {
		return fmt.Errorf("error opening given path: %w", err)
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return fmt.Errorf("error opening given file: %w", err)
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting stats of given file: %w", err)
	}

	fileType := fp.identifyFileType(info)

	if err = fp.handleFile(absolutePath, fileType); err != nil {
		return fmt.Errorf("error handling file: %w", err)
	}

	return nil
}

func (fp *FileProcessor) identifyFileType(info fs.FileInfo) string {

	if info.IsDir() {
		return "directory"
	}

	if strings.LastIndex(info.Name(), ".") > -1 {
		return strings.Split(info.Name(), ".")[1]
	}

	return ""
}

func (fp *FileProcessor) handleFile(path string, fileType string) error {
	handlers := map[string]func(string) error{
		"directory": fp.handleDirectory,
		"docx":      fp.DocxConverter.ConvertDocxFile,
		"zip":       fp.handleZip,
	}

	handler, ok := handlers[fileType]
	if !ok {
		fmt.Printf("Skipping not supported file type: %v\n", path)
		return nil
	}

	return handler(path)
}

func (fp *FileProcessor) handleZip(path string) error {
	dist := strings.Replace(path, ".", "", 1)
	if err := fp.ZipManager.UnZipFile(path, dist); err != nil {
		return err
	}

	if err := fp.handleDirectory(dist); err != nil {
		return err
	}

	if err := fp.ZipManager.ZipDir(dist); err != nil {
		return err
	}

	if err := os.RemoveAll(dist); err != nil {
		return fmt.Errorf("error removing unzipped directory: %w", err)
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("error removing source directory: %w", err)
	}

	return nil
}

func (fp *FileProcessor) handleDirectory(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening directory: %w", err)
	}
	defer dir.Close()

	entries, err := dir.ReadDir(0)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		fp.Process(path + "/" + entry.Name())
	}

	return nil
}
