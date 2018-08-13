package helper // https://github.com/Crosse/font-install/blob/master/font.go

import (
	"bytes"
	"fmt"
	"path"

	"github.com/ConradIrwin/font/sfnt"
)

// FontData describes a font file and the various metadata associated with it.
type FontData struct { // https://github.com/ConradIrwin/font/blob/d797009a8098ca7f6c36a29d0c132a3d39bc4212/sfnt/table_name.go#L77
	Name     string
	Family   string
	FileName string
	Metadata map[sfnt.NameID]string
	Data     []byte
}

// FontExtensions is a list of file extensions that denote fonts.
// Only files ending with these extensions will be installed.
var FontExtensions = map[string]bool{
	".otf": true,
	".ttf": true,
}

// NewFontData creates a new FontData struct.
// fileName is the font's file name, and data is a byte slice containing the font file data.
// It returns a FontData struct describing the font, or an error.
func NewFontData(fileName string, data []byte) (fontData *FontData, err error) {
	if _, ok := FontExtensions[path.Ext(fileName)]; !ok {
		return nil, fmt.Errorf("Not a font: %v", fileName)
	}

	fontData = &FontData{
		FileName: fileName,
		Metadata: make(map[sfnt.NameID]string),
		Data:     data,
	}

	font, err := sfnt.Parse(bytes.NewReader(fontData.Data))
	if err != nil {
		return nil, err
	}

	if font.HasTable(sfnt.TagName) == false {
		return nil, fmt.Errorf("Font %v has no name table", fileName)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		return nil, err
	}

	for _, nameEntry := range nameTable.List() {
		fontData.Metadata[nameEntry.NameID] = nameEntry.String()
	}
	fontData.Name = fontData.Metadata[sfnt.NameFull]
	fontData.Family = fontData.Metadata[sfnt.NamePreferredFamily]
	if fontData.Family == "" {
		if v, ok := fontData.Metadata[sfnt.NameFontFamily]; ok {
			fontData.Family = v
		} else {
			fmt.Println("Font %v has no font family!", fontData.Name)
			// log.Errorf("Font %v has no font family!", fontData.Name)
		}
	}

	if fontData.Name == "" {
		fmt.Println("Font %v has no name! Using file name instead.", fileName)
		// log.Errorf("Font %v has no name! Using file name instead.", fileName)
		fontData.Name = fileName
	}

	return
}
