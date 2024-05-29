package tar_utils

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
)

//goland:noinspection GoUnhandledErrorResult
func Path2TarReader(sourcePath string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	err := filepath.Walk(sourcePath, func(currentPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Walk in error: %v", err)
		}
		// 创建 tar 包头部信息
		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			log.Fatalf("FileInfoHeader error: %v", err)
		}
		relPath, err := filepath.Rel(sourcePath, currentPath)
		if err != nil {
			log.Fatalf("Rel error: %v", err)
		}
		header.Name = relPath
		if err := tw.WriteHeader(header); err != nil {
			log.Fatalf("WriteHeader error: %v", err)
		}
		// 跳过非常规文件
		if !fileInfo.Mode().IsRegular() {
			return nil
		}
		currentFile, err := os.Open(currentPath)
		if err != nil {
			log.Fatalf("Open error: %v", err)
		}
		defer currentFile.Close()
		// 创建 tar 包文件数据
		if _, err := io.Copy(tw, currentFile); err != nil {
			log.Fatalf("Copy error: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Walk error: %v", err)
	}
	return buf, nil
}
