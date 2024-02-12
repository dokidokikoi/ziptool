package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ziptool/flag"
	"ziptool/log"
	"ziptool/plugins"
	"ziptool/plugins/rar"
	"ziptool/plugins/sevenz"
	"ziptool/plugins/ziptool"
	"ziptool/woker"

	"github.com/h2non/filetype"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	copyLimit int64 = 100 << 20
)

const (
	zip      = "application/zip"
	sevenzip = "application/x-7z-compressed"
	gzip     = "application/gzip"
	rarzip   = "application/vnd.rar"
)

var m = map[string]plugins.IPlugin{
	zip:      &ziptool.ZipTool{},
	sevenzip: &sevenz.Sevenz{},
	rarzip:   &rar.Rar{},
}

var rootCmd = &cobra.Command{
	Use:   "ziptool",
	Short: "bulk compress or extract",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		src, ok := viper.Get("Src").(string)
		if src == "" || !ok {
			return
		}
		dest, ok := viper.Get("Dest").(string)
		if dest == "" || !ok {
			dest = filepath.Join(src, "ziptooltmp")
		}
		passwords := viper.GetStringSlice("Passwords")
		err := os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			panic(err)
		}
		filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
			if path == dest {
				return errors.New("skip tmp dir")
			}
			if info.IsDir() {
				return err
			}
			rel, err := filepath.Rel(src, path)
			if err != nil {
				rel = path
			}
			target := filepath.Join(dest, rel)
			targetTmp := strings.Split(target, "/")
			var targetDir string
			if len(targetTmp) > 1 {
				targetDir = strings.Join(targetTmp[:len(targetTmp)-1], "/")
				os.MkdirAll(targetDir, os.ModePerm)
			} else {
				return nil
			}
			if target == path {
				return nil
			}
			p := dis(path)
			if p == nil {
				f, err := os.Open(path)
				if err != nil {
					log.Add(log.Log{Path: path, Result: err.Error(), CreateTime: time.Now(), Flag: log.FAIL})
					return nil
				}
				buf := make([]byte, 1024)
				io.ReadFull(f, buf)
				var ty = log.ETC
				if filetype.IsImage(buf) {
					ty = log.IMAGE
				}
				info, _ := f.Stat()
				if info.Size() < copyLimit {
					woker.AddTask(woker.Task{
						Src:  path,
						Dest: target,
						Type: ty,
					})
				}
			} else {
				if p.ExtractAble(path) {
					woker.AddTask(woker.Task{
						Src:       path,
						Dest:      target,
						TmpDir:    filepath.Join(targetDir, fmt.Sprintf("%d", time.Now().Unix())),
						Zipper:    p,
						Passwords: passwords,
						Type:      log.ARCHIVE,
					})
				}
			}

			return nil
		})
		woker.RunBar()
		go func() {
			c := make(chan os.Signal)
			signal.Notify(c, syscall.SIGINT)
			<-c
			fmt.Println()
			fmt.Println("program exiting ...")
			woker.Close()
		}()
		flag.Wait.Wait()
		log.Close()
		flag.LogWait.Wait()
	},
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "/data/yaml/system-api.yaml", "config file")
}

func dis(file string) plugins.IPlugin {
	buf, _ := os.ReadFile(file)

	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		return nil
	}

	switch kind.MIME.Value {
	case zip:
		return m[zip]
	case sevenzip:
		return m[sevenzip]
	case gzip:
		return m[gzip]
	case rarzip:
		return m[rarzip]
	}

	return nil
}
