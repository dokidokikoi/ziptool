package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/h2non/filetype"
)

func TestFilePath(t *testing.T) {
	filepath.Walk("..", func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path, info.Name())
		if info.Name() == "cmd" {
			return errors.New("no more")
		}
		return err
	})
}

func TestPathJoin(t *testing.T) {
	fmt.Println(filepath.Join("/loacl/test", "./image/test.png"))
}

func TestDirFiles(t *testing.T) {
	dirs, _ := os.ReadDir("..")
	for _, d := range dirs {
		fmt.Println(d.Name())
	}
}

func TestType(t *testing.T) {
	ty, err := filetype.MatchFile("../plugins/example/image.7z.002")
	if err != nil {
		panic(err)
	}

	fmt.Println(ty.MIME.Value)
}

func TestRegexp(t *testing.T) {
	ok, err := regexp.MatchString(`^[\w`+"\u4e00-\u9fa5"+`]+\.part\d+\.rar$`, "12sda.part1.rar")
	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
}
