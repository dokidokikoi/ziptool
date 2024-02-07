package rar

import (
	"log"
	"os"
	"testing"
)

var rarS = &Rar{}

func TestCompress(t *testing.T) {
	err := rarS.Compress([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.rar")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompressWithPwd(t *testing.T) {
	err := rarS.CompressWithPwd([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.rar", "123")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompressWithPwdVolumes(t *testing.T) {
	err := rarS.CompressWithPwdVolumes([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.rar", "123", "1000k")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompressVolumes(t *testing.T) {
	err := rarS.CompressVolumes([]string{"../example/sen.jpeg", "../example/wallhaven-1p2mxw.jpg"}, "../example/image.rar", "2m")
	if err != nil {
		log.Fatal(err)
	}
}

func TestExtract(t *testing.T) {
	os.Mkdir("1", os.ModePerm)
	err := rarS.Extract("../example/image.part4.rar", "./1", "")
	if err != nil {
		log.Fatal(err)
	}
}
