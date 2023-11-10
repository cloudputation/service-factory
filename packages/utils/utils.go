package utils

import (
	"os"
	"log"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)


func CreateDir(dirName string) error {
	fmt.Printf("Checking if directory %s already exists.\n", dirName)
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		// Check if the error is because the directory already exists
		if os.IsExist(err) {
			fmt.Printf("Directory %s already exists.\n", dirName)
		} else {
			// Handle different kinds of errors here
			log.Fatalf("Error creating directory: %v", err)
		}
	} else {
		fmt.Printf("Directory %s created.\n", dirName)
	}

	return nil
}

func DeleteDir(dirName string) error {
	err := os.RemoveAll(dirName)
	if err != nil {
		log.Fatalf("Error removing directory: %v", err)
	}
	fmt.Println("Directory removed.")

	return nil
}

func CopyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(srcPath string, dstPath string) error {
	// Open source file for reading
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Get the source file's permissions
	srcInfo, err := src.Stat()
	if err != nil {
		return err
	}

	// Create destination file with the source file's permissions
	dst, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the contents of the source file into the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	fmt.Printf("File %s has been copied to %s\n", srcPath, dstPath)

	return nil
}
