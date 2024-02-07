package log

import (
	"encoding/json"
	"os"
	"sync"
	"time"
	"ziptool/flag"
)

const (
	SUCC = "success"
	FAIL = "fail"
)

const (
	ETC     = "etc"
	IMAGE   = "image"
	ARCHIVE = "archive"
)

var closeOnce sync.Once

type Log struct {
	Path       string    `json:"path"`
	Result     string    `json:"result"`
	CreateTime time.Time `json:"-"`
	CreateStr  string    `json:"create_time"`
	Flag       string    `json:"flag"`
	Type       string    `json:"type"`
}

var lc LogCtl

type LogCtl struct {
	c       chan Log
	close   chan struct{}
	success []Log
	failed  []Log
}

func init() {
	lc = LogCtl{
		c:     make(chan Log, 30),
		close: make(chan struct{}),
	}
	logfile := time.Now().Format("2006-01-02_15:04:05")
	flag.Wait.Add(1)
	go func() {
		defer flag.Wait.Done()
		for {
			select {
			case <-lc.close:
				if len(lc.success) > 0 {
					lsf, err := os.Create("log_" + logfile + "_success.json")
					if err != nil {
						panic(err)
					}
					e := json.NewEncoder(lsf)
					e.Encode(lc.success)
					lsf.Close()
				}
				if len(lc.failed) > 0 {
					lff, err := os.Create("log_" + logfile + "_failed.json")
					if err != nil {
						panic(err)
					}
					e := json.NewEncoder(lff)
					e.Encode(lc.failed)
					lff.Close()
				}

				return
			case l := <-lc.c:
				l.CreateStr = l.CreateTime.Format("2006-01-02 15:04:05")
				if l.Flag == SUCC {
					lc.success = append(lc.success, l)
				} else {
					lc.failed = append(lc.failed, l)
				}
			}
		}
	}()
}

func Add(l Log) {
	lc.c <- l
}

func Success(l Log) {
	l.Flag = SUCC
	lc.c <- l
}

func Failed(l Log) {
	l.Flag = FAIL
	lc.c <- l
}

func Close() {
	closeOnce.Do(func() {
		close(lc.close)
	})
}
