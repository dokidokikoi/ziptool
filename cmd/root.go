package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
	cfgFile string
	bufSize int64 = 1024 * 1024
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
				buf, _ := os.ReadFile(path)
				var ty = log.ETC
				if filetype.IsImage(buf) {
					ty = log.IMAGE
				}
				f, err := os.Open(path)
				if err != nil {
					log.Failed(log.Log{Path: path, Result: err.Error(), CreateTime: time.Now(), Type: ty})
					return nil
				}
				targetF, err := os.Create(target)
				if err != nil {
					log.Failed(log.Log{Path: path, Result: err.Error(), CreateTime: time.Now(), Type: ty})
					return nil
				}
				for {
					_, err := io.CopyN(targetF, f, bufSize)
					if err != nil {
						if errors.Is(err, io.EOF) {
							break
						} else {
							log.Add(log.Log{Path: path, Result: err.Error(), CreateTime: time.Now(), Type: ty})
							return nil
						}
					}
				}
				log.Success(log.Log{Path: path, Result: "success", CreateTime: time.Now(), Type: ty})
			} else {
				if p.ExtractAble(path) {
					woker.AddTask(woker.Task{
						Src:       path,
						Dest:      target,
						TmpDir:    filepath.Join(targetDir, fmt.Sprintf("%d", time.Now().Unix())),
						Zipper:    p,
						Passwords: passwords,
					})
				}
			}

			return nil
		})
		woker.RunBar()

		flag.Wait.Wait()
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
