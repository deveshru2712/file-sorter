package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func SortFile(filePath string) {
	data, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// in case folder does not exists
		fmt.Printf("%s does not exists\n", filePath)
		return
	} else if err != nil {
		// in case of some other error
		fmt.Printf("Error accessing %s\n", filePath)
		return
	}

	if data.IsDir() {
		isEmpty, err := isDirEmpty(filePath)

		if isEmpty {
			fmt.Printf("Folder is empty")
			return
		}
		if err != nil {
			fmt.Printf("Error occurred while reading the dir %s", filePath)
			return
		}

		// print the folder that i need to create
		res := foldersToCreate(filePath)
		for _, ext := range res {
			// creating those folder the same dir
			newFilePath := filepath.Join(filePath, ext)
			err := os.Mkdir(newFilePath, 0755)
			if err != nil {
				fmt.Println("Error creating folder:", err)
				return
			}
		}

		fmt.Println("folders created")

	} else {
		// if path is for a file
		fmt.Println("soon")
	}
}

func isDirEmpty(filePath string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	// closing the file
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func foldersToCreate(filePath string) []string {
	extMap := make(map[string]bool)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil
	}

	for _, e := range files {
		if e.IsDir() {
			continue
		}

		ext := filepath.Ext(e.Name())
		if ext != "" {
			folderName := getCategory(ext)
			extMap[folderName] = true
		}
	}

	var res []string
	for ext := range extMap {
		// removing the dot
		res = append(res, ext)
	}

	return res
}

func getCategory(extension string) string {
	switch extension {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "images"

	case ".pdf", ".doc", ".docx", ".txt":
		return "documents"

	case ".mp4", ".mkv", ".avi":
		return "videos"

	case ".mp3", ".wav":
		return "audio"

	default:
		return extension
	}
}
