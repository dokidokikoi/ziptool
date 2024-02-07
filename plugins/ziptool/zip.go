package ziptool

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mzky/zip"
)

type ZipTool struct {
}

func (zt ZipTool) Compress(srcPath []string, destPath string) error {
	return Zip(destPath, "", srcPath)
}
func (zt ZipTool) CompressWithPwd(srcPath []string, destPath, password string) error {
	return Zip(destPath, password, srcPath)
}
func (zt ZipTool) CompressVolumes(srcPath []string, destPath string, size string) error {
	fmt.Println("not implement")
	return nil
}
func (zt ZipTool) CompressWithPwdVolumes(srcPath []string, destPath, password string, size string) error {
	fmt.Println("not implement")
	return nil
}

func (zt ZipTool) Extract(srcPath, destPath, password string) error {
	return UnZip(srcPath, password, destPath)
}

func (zt ZipTool) ExtractAble(srcPath string) bool {
	return true
}

func (zt ZipTool) FileName(srcPath string) string {
	tmp := strings.Split(srcPath, ".")
	return strings.Join(tmp[:len(tmp)-1], ".")
}

// password值可以为空""
func Zip(zipPath, password string, fileList []string) error {
	if len(fileList) < 1 {
		return fmt.Errorf("将要压缩的文件列表不能为空")
	}
	fz, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(fz)
	defer zw.Close()

	for _, fileName := range fileList {
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}

		// 写入文件的头信息
		var w io.Writer
		if password != "" {
			w, err = zw.Encrypt(fileName, password, zip.AES256Encryption)
		} else {
			w, err = zw.Create(fileName)
		}

		if err != nil {
			return err
		}

		// 写入文件内容
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return zw.Flush()
}

// password值可以为空""
// 当decompressPath值为"./"时，解压到相对路径
func UnZip(zipPath, password, decompressPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if password != "" {
			if f.IsEncrypted() {
				f.SetPassword(password)
			} else {
				return errors.New("must be encrypted")
			}
		}
		fp := filepath.Join(decompressPath, f.Name)
		dir, _ := filepath.Split(fp)
		if dir != "" {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		w, err := os.Create(fp)
		if nil != err {
			return err
		}

		fr, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
		w.Close()
	}
	return nil
}
