package rar

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	rarCmd     = "rar"
	addCmd     = "a"
	extractCmd = "e"

	volumesSw       = "-v{size}" // [k,b] size=<size>*1000 [*1024, *1]
	sizePattern     = "{size}"
	passwordSw      = "-p{password}"
	passwordPattern = "{password}"
	assumeYSw       = "-y"
)

type Rar struct {
}

func (r Rar) Compress(srcPath []string, destPath string) error {
	return r.CompressWithPwd(srcPath, destPath, "")
}
func (r Rar) CompressWithPwd(srcPath []string, destPath, password string) error {
	var args = []string{
		addCmd, destPath,
	}
	args = append(args, srcPath...)
	if password != "" {
		args = append(args, strings.ReplaceAll(passwordSw, passwordPattern, password))
	}
	cmd := exec.Command(rarCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Start()
	err := cmd.Wait()
	return err
}
func (r Rar) CompressVolumes(srcPath []string, destPath string, size string) error {
	return r.CompressWithPwdVolumes(srcPath, destPath, "", size)
}
func (r Rar) CompressWithPwdVolumes(srcPath []string, destPath, password string, size string) error {
	var args = []string{
		addCmd, destPath, strings.ReplaceAll(volumesSw, sizePattern, size),
	}
	args = append(args, srcPath...)
	if password != "" {
		args = append(args, strings.ReplaceAll(passwordSw, passwordPattern, password))
	}
	cmd := exec.Command(rarCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Start()
	err := cmd.Wait()
	return err
}

func (r Rar) Extract(srcPath, destPath, password string) error {
	var args = []string{
		extractCmd, srcPath, destPath, assumeYSw,
	}
	if password != "" {
		args = append(args, strings.ReplaceAll(passwordSw, passwordPattern, password))
	}
	cmd := exec.Command(rarCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Start()
	err := cmd.Wait()
	return err
}

func (r Rar) ExtractAble(srcPath string) bool {
	tmp := strings.Split(srcPath, "/")
	filename := tmp[len(tmp)-1]
	if strings.HasSuffix(filename, ".part1.rar") {
		return true
	}
	ok, err := regexp.MatchString(`^[\w`+"\u4e00-\u9fa5"+`]+\.part\d+\.rar$`, filename)
	if err != nil {
		return false
	}
	if ok {
		return false
	}
	if strings.HasSuffix(filename, ".rar") {
		return true
	}
	return false
}

func (r Rar) FileName(srcPath string) string {
	tmp := strings.Split(srcPath, "/")
	filename := tmp[len(tmp)-1]
	tmp = strings.Split(srcPath, ".")
	ok, err := regexp.MatchString(`^[\w`+"\u4e00-\u9fa5"+`]+\.part\d+\.rar$`, filename)
	if err != nil {
		return ""
	}
	if ok {
		return strings.Join(tmp[:len(tmp)-2], ".")
	}
	return strings.Join(tmp[:len(tmp)-1], ".")
}
