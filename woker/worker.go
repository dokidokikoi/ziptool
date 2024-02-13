package woker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	"ziptool/flag"
	"ziptool/log"
	"ziptool/plugins"
	"ziptool/plugins/sevenz"
)

var bufSize int64 = 1024 * 1024
var wpIns = wokerPool{}
var sizeLimit int64 = 3 << 30
var volumesSize = "2000m"
var (
	closeOnce  sync.Once
	finishOnce sync.Once
)

type wokerPool struct {
	pool  chan struct{}
	tasks chan Task
	close chan struct{}

	poolSize int
	finished chan struct{}
}

type Task struct {
	Src       string
	Dest      string
	TmpDir    string
	Passwords []string
	Size      int64
	Zipper    plugins.IPlugin
	Type      string
}

func (t Task) Run() error {
	if t.Type == log.ARCHIVE {
		err := os.MkdirAll(t.TmpDir, os.ModePerm)
		if err != nil {
			return err
		}
		defer os.RemoveAll(t.TmpDir)

		err = t.Zipper.Extract(t.Src, t.TmpDir, "")
		for i := 0; err != nil && i < len(t.Passwords); i++ {
			err = t.Zipper.Extract(t.Src, t.TmpDir, t.Passwords[i])
		}
		if err != nil {
			return err
		}
		fmt.Println("Extract completed", t.Src)
		dest := t.Zipper.FileName(t.Dest) + ".7z"
		sz := sevenz.Sevenz{}
		files := []string{}
		dirs, err := os.ReadDir(t.TmpDir)
		if err != nil {
			return err
		}
		for _, d := range dirs {
			files = append(files, filepath.Join(t.TmpDir, d.Name()))
		}
		if dirSize(t.TmpDir) > sizeLimit {
			err = sz.CompressVolumes(files, dest, volumesSize)
			if err != nil {
				return err
			}
		} else {
			err = sz.Compress(files, dest)
			if err != nil {
				return err
			}
		}
	} else {
		f, err := os.Open(t.Src)
		if err != nil {
			return err
		}
		targetF, err := os.Create(t.Dest)
		if err != nil {
			return err
		}
		for {
			_, err := io.CopyN(targetF, f, bufSize)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					return err
				}
			}
		}
	}

	return nil
}

func init() {
	wpIns = wokerPool{
		pool:     make(chan struct{}, 10),
		poolSize: 10,
		tasks:    make(chan Task, 20),
		close:    make(chan struct{}),
		finished: make(chan struct{}),
	}
	Run()
}

func Run() {
	for i := 0; i < wpIns.poolSize; i++ {
		flag.Wait.Add(1)
		go func() {
			defer flag.Wait.Done()
			for {
				select {
				case <-wpIns.close:
					return
				case t := <-wpIns.tasks:
					err := t.Run()
					if err != nil {
						log.Add(log.Log{Path: t.Src, Result: err.Error(), CreateTime: time.Now(), Flag: log.FAIL, Type: t.Type})
					} else {
						log.Add(log.Log{Path: t.Src, Result: "success", CreateTime: time.Now(), Flag: log.SUCC, Type: t.Type})
					}
					Process(t)
				case <-wpIns.finished:
					if len(wpIns.tasks) == 0 {
						Close()
					}
				}
			}
		}()
	}
}

func Close() {
	closeOnce.Do(func() {
		close(wpIns.close)
	})
}

func Finished() {
	finishOnce.Do(func() {
		close(wpIns.finished)
	})
}

func Add(t Task) {
	wpIns.tasks <- t
}

func dirSize(path string) int64 {
	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return size
}
