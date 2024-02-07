package sevenz

import (
	"log"
	"os"
	"testing"
)

var sevenZS = &Sevenz{}

func TestCompress(t *testing.T) {
	err := sevenZS.Compress([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.7z")
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
	err := sevenZS.Extract("../example/image.7z.003", "./1", "")
	if err != nil {
		log.Fatal(err)
	}
}
