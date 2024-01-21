package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

func main() {

	if len(os.Args) > 2 || len(os.Args) == 1 {
		fmt.Println("Please provide only one argument: path to file or directory")
		return
	}

	arg := os.Args[1]

	path, err := filepath.Abs(arg)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening given path: %w", err))
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening given path: %w", err))
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(fmt.Errorf("error opening given path: %w", err))
	}

	switch {
	case info.IsDir():
		if err := handleDirectory(path); err != nil {
			log.Fatal(err)
		}

	case strings.HasSuffix(info.Name(), ".docx"):
		if err := convertDocxFile(path); err != nil {
			log.Fatal(err)
		}

	case strings.HasSuffix(info.Name(), ".zip"):
		dist := strings.Replace(path, ".", "", 1)
		if err := unZipFile(path, dist); err != nil {
			log.Fatal(err)
		}

		if err := handleDirectory(dist); err != nil {
			log.Fatal(err)
		}

		if err := zipDir(dist); err != nil {
			log.Fatal(err)
		}

		if err := os.RemoveAll(dist); err != nil {
			log.Fatal(fmt.Errorf("error removing old directory: %w", err))
		}

		if err := os.RemoveAll(path); err != nil {
			log.Fatal(fmt.Errorf("error removing old directory: %w", err))
		}

	default:
		fmt.Printf("Skipping not supported file type: %v\n", info.Name())
	}
}

func handleDirectory(path string) error {
	fmt.Printf("Opening directory: %v\n", path)
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
		info, _ := entry.Info()
		switch {
		case info.IsDir():
			if err := handleDirectory(path + "/" + info.Name()); err != nil {
				fmt.Println(err)
			}
		case strings.HasSuffix(info.Name(), ".docx"):
			if err := convertDocxFile(path + "/" + info.Name()); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("Skipping not supported file type: %v\n", info.Name())
		}
	}

	return nil
}

func unZipFile(path string, distPath string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("error unzipping file: %w", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(distPath, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(distPath)+string(os.PathSeparator)) {
			return fmt.Errorf("error unzipping file: %w", err)
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("error unzipping file: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("error unzipping file: %w", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("error unzipping file: %w", err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("error unzipping file: %w", err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("error unzipping file: %w", err)
		}

		defer dstFile.Close()
		defer fileInArchive.Close()
	}
	return nil
}

func zipDir(path string) error {
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error zipping directory: %w", err)
	}

	zipFileName := path + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("error zipping directory: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error zipping directory: %w", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error zipping directory: %w", err)
		}
		defer file.Close()

		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return fmt.Errorf("error zipping directory: %w", err)
		}

		if !fileInfo.IsDir() {
			zipEntry, err := zipWriter.Create(relPath)
			if err != nil {
				return fmt.Errorf("error zipping directory: %w", err)
			}

			_, err = io.Copy(zipEntry, file)
			if err != nil {
				return fmt.Errorf("error zipping directory: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func convertDocxFile(path string) error {
	r, err := docx.ReadDocxFile(path)
	if err != nil {
		return fmt.Errorf("error reading docx file: %w", err)
	}
	defer r.Close()

	file := r.Editable()
	content := file.GetContent()

	// Replace dividing line with dashes
	dividingLine := `<w:r w:rsidDel="00000000" w:rsidR="00000000" w:rsidRPr="00000000"><w:pict><v:rect style="width:0.0pt;height:1.5pt" o:hr="t" o:hrstd="t" o:hralign="center" fillcolor="#A0A0A0" stroked="f"/></w:pict></w:r>`
	dashes := `<w:r w:rsidDel="00000000" w:rsidR="00000000" w:rsidRPr="00000000"><w:rPr><w:rtl w:val="0"/></w:rPr><w:t xml:space="preserve">---</w:t></w:r>`
	content = strings.ReplaceAll(content, dividingLine, dashes)

	// Remove hyperlinks underlines
	arr := strings.Split(content, "<w:hyperlink")
	for i, element := range arr {
		// make sure to only replace the occurence of the underline tag inside the hyperlink tag not anywhere else
		internalArr := strings.Split(element, "</w:hyperlink>")
		internalArr[0] = strings.ReplaceAll(internalArr[0], `<w:u w:val="single"/>`, "")
		arr[i] = strings.Join(internalArr, "</w:hyperlink>")
	}
	content = strings.Join(arr, "<w:hyperlink")

	// save the file with the new content
	file.SetContent(content)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error removing old docx file: %w", err)
	}
	// write previous content back to file if error occurs
	if err := file.WriteToFile(path); err != nil {
		file.SetContent(r.Editable().GetContent())
		file.WriteToFile(path)
		return fmt.Errorf("error creating new docx file: %w", err)
	}
	fmt.Printf("Successfully converted: %v\n", path)
	return nil
}
