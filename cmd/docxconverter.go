package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

type DocxConverter struct{}

const (
	// DividingLine is the string that represents the dividing line in a docx file
	DividingLine = `<w:r w:rsidDel="00000000" w:rsidR="00000000" w:rsidRPr="00000000"><w:pict><v:rect style="width:0.0pt;height:1.5pt" o:hr="t" o:hrstd="t" o:hralign="center" fillcolor="#A0A0A0" stroked="f"/></w:pict></w:r>`
	// Dashes is the string that represents the three dashes in a docx file
	Dashes = `<w:r w:rsidDel="00000000" w:rsidR="00000000" w:rsidRPr="00000000"><w:rPr><w:rtl w:val="0"/></w:rPr><w:t xml:space="preserve">---</w:t></w:r>`
	// HyperlinkUnderline is the string that represents the underline tag in a hyperlink tag in a docx file
	HyperlinkUnderline  = `<w:u w:val="single"/>`
	HyperlinkOpeningTag = "<w:hyperlink>"
	HyperlinkClosingTag = "</w:hyperlink>"
)

func (c *DocxConverter) ConvertDocxFile(path string) error {
	r, err := docx.ReadDocxFile(path)
	if err != nil {
		return fmt.Errorf("error reading docx file: %w", err)
	}
	defer r.Close()

	file := r.Editable()
	content := file.GetContent()

	content = c.replaceDividingLine(content)
	content = c.removeHyperlinksUnderlines(content)

	// save the file with the new content
	file.SetContent(content)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error removing path: %w", err)
	}
	// write previous content back to file if error occurs
	if err := file.WriteToFile(path); err != nil {
		file.SetContent(r.Editable().GetContent())
		file.WriteToFile(path)
		return fmt.Errorf("error writing new content to docx file: %w", err)
	}
	fmt.Printf("Successfully converted: %v\n", path)
	return nil
}

func (c *DocxConverter) replaceDividingLine(content string) string {
	return strings.ReplaceAll(content, DividingLine, Dashes)
}

func (c *DocxConverter) removeHyperlinksUnderlines(content string) string {
	// cut last character of the opening tag to make sure to remove the underline tag inside the hyperlink tag
	arr := strings.Split(content, HyperlinkOpeningTag[:len(HyperlinkOpeningTag)-1])
	for i, element := range arr {
		// make sure to only replace the occurrence of the underline tag inside the hyperlink tag not anywhere else
		internalArr := strings.Split(element, HyperlinkClosingTag)
		internalArr[0] = strings.ReplaceAll(internalArr[0], HyperlinkUnderline, "")
		arr[i] = strings.Join(internalArr, HyperlinkClosingTag)
	}
	return strings.Join(arr, HyperlinkOpeningTag[:len(HyperlinkOpeningTag)-1])
}
