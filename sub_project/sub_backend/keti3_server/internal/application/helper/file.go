package helper

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// DownloadFile 本地文件路径
func DownloadFile(fileURL, module string) (string, error) {
	var (
		err                            error
		localPath, localFile, filename string
	)
	if localPath, err = MakeTempDir(fmt.Sprintf("import-file/%s", module)); err != nil {
		return "", err
	}

	filename = path.Base(fileURL)
	if filename == "" {
		filename = fmt.Sprintf("%d.xlsx", time.Now().UnixMilli())
	}
	localFile = fmt.Sprintf("%s/%s", localPath, filename)
	if err = downloadFile(localFile, fileURL); err != nil {
		return "", err
	}
	return localFile, nil
}

// MakeTempDir 生成临时文件夹
func MakeTempDir(name string, paths ...string) (string, error) {
	var (
		dir string
		err error
	)
	dir = fmt.Sprintf("./logs/%s", name)
	if err = os.MkdirAll(dir, 0750); err != nil {
		os.RemoveAll(dir)
		return "", err
	}
	if dir, err = os.MkdirTemp(dir, ""); err != nil {
		return "", err
	}
	if len(paths) > 0 {
		for _, path := range paths {
			dir = filepath.Join(dir, path)
		}
	}
	if err = os.MkdirAll(dir, 0750); err != nil {
		os.RemoveAll(dir)
		return "", err
	}
	return dir, nil
}

// DownloadFile 下载文件
func downloadFile(localFile, url string) error {
	var (
		resp *http.Response
		f    *os.File
		err  error
	)

	if resp, err = http.Get(url); err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = os.MkdirAll(filepath.Dir(localFile), 0775); err != nil {
		return err
	}
	if f, err = os.Create(localFile); err != nil {
		return err
	}
	if _, err = io.Copy(f, resp.Body); err != nil {
		return err
	}
	return nil
}

// ZipDirectory 将指定目录及其内容压缩到目标zip文件中
func ZipDirectory(sourceDir, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径以去掉源目录前缀
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// 忽略根目录本身
		if relPath == "." {
			return nil
		}

		// 创建压缩文件头信息
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件则写入数据
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
