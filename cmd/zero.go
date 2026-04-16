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

			// check if folder already exists
			folderPath := filepath.Join(filePath, ext)
			if _, err := os.Stat(folderPath); err == nil {
				// folder already exists, skip creating it
				continue
			} else if !os.IsNotExist(err) {
				// some other error occurred while checking the folder
				fmt.Printf("Error checking folder %s: %v\n", folderPath, err)
				return
			}

			// creating those folder the same dir
			newFilePath := filepath.Join(filePath, ext)
			err := os.Mkdir(newFilePath, 0755)
			if err != nil {
				fmt.Println("Error creating folder:", err)
				return
			}
		}

		// moving the files to those folders
		files, err := os.ReadDir(filePath)
		if err != nil {
			fmt.Printf("Error occurred while reading the dir %s", filePath)
			return
		}

		for _, f := range files {
			// if it is a folder then skip it
			if f.IsDir() {
				continue
			}

			src := filepath.Join(filePath, f.Name())
			dist, err := getDestination(src)
			if err != nil {
				fmt.Printf("Error getting destination for file %s: %v\n", src, err)
				continue
			}
			err = moveFile(src, dist)
			if err != nil {
				fmt.Printf("Error moving file %s to %s: %v\n", src, dist, err)
				continue
			}
			fmt.Printf("Moved file %s to %s\n", src, dist)
		}
		fmt.Println("Sorting completed!")

	} else {
		// if path is for a file
		fmt.Println("soon")
	}
}

func getDestination(filePath string) (string, error) {
	ext := filepath.Ext(filePath)
	if ext == "" {
		return "", fmt.Errorf("file has no extension")
	}

	folderName := getCategory(ext)
	dist := filepath.Join(filepath.Dir(filePath), folderName, filepath.Base(filePath))
	return dist, nil
}

func moveFile(src, dist string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dist, input, 0644)
	if err != nil {
		return err
	}
	return os.Remove(src)
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
