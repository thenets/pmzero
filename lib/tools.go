package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

// HasLettersOnly returns true if all string is letters
func HasLettersOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// CopyFile from src to dst
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// IsDirectory returns true if is a directory
func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// LoadConfigFiles loads a file or directory to
func LoadConfigFiles(filePath string) error {
	var err error

	// Check if is a dir
	if IsDirectory(filePath) {
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			var absoluteFilePath, _ = filepath.Abs(filePath + "/" + f.Name())
			LoadConfigFiles(absoluteFilePath)
		}
		return err
	}

	// Check if is a deployment
	var data = GetDeploymentByFilePath(filePath)
	if data.Type == "deployment" {
		LoadDeploymentFile(filePath)
	} else {
		log.Fatalf("[ERROR] config file type not supported: %v\n", data)
	}

	return err
}
