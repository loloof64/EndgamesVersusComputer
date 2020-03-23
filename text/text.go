package text

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wailsapp/wails"
)

// TextFileManager is an utility for reading and saving text files
type TextFileManager struct {
	Runtime *wails.Runtime
}

// NewTextFileManager returns a new text loader.
func NewTextFileManager() *TextFileManager {
	return &TextFileManager{}
}

// WailsInit initialize the structure for Wails
func (loader *TextFileManager) WailsInit(r *wails.Runtime) error {
	loader.Runtime = r
	return nil
}

// GetTextFileContent let user select a text file and returns its content
// [output] string: the text content if no error
// otherwise the error code:
// #ErrorReadingFile: the file could not be read
func (loader *TextFileManager) GetTextFileContent() string {
	selectedFile := loader.Runtime.Dialog.SelectFile()

	data, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		return "#ErrorReadingFile"
	}

	result := string(data)

	return result
}

// GetTextFileContentWithPathProviden let you read content from a text file without file chooser.
// path string: the absolute path
// [output] string: the text content if no error
// otherwise the error code:
// #ErrorReadingFile: the file could not be read
func (loader TextFileManager) GetTextFileContentWithPathProviden(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "#ErrorReadingFile"
	}

	result := string(data)

	return result
}

// SaveTextFile lets user choose a file and save text content into it.
// content string: the content to save
// [output] string: the error if any
// #ErrorSavingFile: the file could not be saved
// #ErrorClosingFile: the new file could not be closed
func (loader TextFileManager) SaveTextFile(content string) string {
	selectedFile := loader.Runtime.Dialog.SelectSaveFile()

	file, err := os.Create(selectedFile)
	if err != nil {
		fmt.Println(err)
		return "#ErrorSavingFile"
	}

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return "#ErrorSavingFile"
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return "#ErrorClosingFile"
	}

	return ""
}

// SaveTextFileWithPathProviden save text content into given file path.
// path string: the absolute path
// content string: the content to save
// [output] string: the error if any
// #ErrorSavingFile: the file could not be saved
// #ErrorClosingFile: the new file could not be closed
func (loader TextFileManager) SaveTextFileWithPathProviden(path string, content string) string {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return "#ErrorSavingFile"
	}

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return "#ErrorSavingFile"
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return "#ErrorClosingFile"
	}

	return ""
}
