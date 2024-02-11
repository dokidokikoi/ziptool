package sevenz

import (
	"log"
	"os"
	"testing"
)

var sevenZS = &Sevenz{}

func TestCompress(t *testing.T) {
	err := sevenZS.Compress([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image1.7z")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompressWithPwd(t *testing.T) {
	err := sevenZS.CompressWithPwd([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.7z", "123")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompressWithPwdVolumes(t *testing.T) {
	err := sevenZS.CompressWithPwdVolumes([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.7z", "123", "1000k")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompressVolumes(t *testing.T) {
	err := sevenZS.CompressVolumes([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.7z", "2m")
	if err != nil {
		log.Fatal(err)
	}
}

func TestExtract(t *testing.T) {
	os.Mkdir("1", os.ModePerm)
	err := sevenZS.Extract("/Users/doki/Desktop/test/files/file/炎孕：异世界超エロ恶魔学园_AI精翻汉化版+存档+DLC【神拔作新汉化全动态】_7z.001", "./1", "")
	if err != nil {
		log.Fatal(err)
	}
}
