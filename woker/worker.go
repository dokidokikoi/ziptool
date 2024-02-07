package woker

import (
	"os"
	"path/filepath"
	"sync"
	"time"
	"ziptool/flag"
	"ziptool/log"
	"ziptool/plugins"
	"ziptool/plugins/sevenz"
)

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
}

func (t Task) Run() error {
	err := os.MkdirAll(t.TmpDir, os.ModePerm)
	if err != nil {
		return err
	}
	defer os.RemoveAll(t.TmpDir)
	f, err := os.Open(t.Src)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}
	err = t.Zipper.Extract(t.Src, t.TmpDir, "")
	for i := 0; err != nil && i < len(t.Passwords); i++ {
		err = t.Zipper.Extract(t.Src, t.TmpDir, t.Passwords[i])
	}
	if err != nil {
		return err
	}
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
	if info.Size() > sizeLimit {
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
					log.Close()
					return
				case t := <-wpIns.tasks:
					err := t.Run()
					if err != nil {
						log.Add(log.Log{Path: t.Src, Result: err.Error(), CreateTime: time.Now(), Flag: log.FAIL, Type: log.ARCHIVE})
					} else {
						log.Add(log.Log{Path: t.Src, Result: "success", CreateTime: time.Now(), Flag: log.SUCC, Type: log.ARCHIVE})
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
