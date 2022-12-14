package service

import (
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	data, _ := io.ReadAll(src)
	createDir(dst)
	return os.WriteFile(dst, data, 0666)
}

func createDir(dst string) {
	dst = dst[0:strings.LastIndex(dst, "/")]
	if _, err := os.Stat(dst); err != nil {
		if os.IsNotExist(err) {
			//创建文件夹
			_ = os.MkdirAll(dst, os.ModeDir)
		}
	}
}
