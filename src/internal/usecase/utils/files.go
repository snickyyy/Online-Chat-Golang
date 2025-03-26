package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadFile(fileHeader *multipart.FileHeader, path string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

func DeleteIfExist(path string) bool {
	return os.Remove(path) == nil
}
