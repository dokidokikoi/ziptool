package sevenz

import "strings"

const sevenCmd = "7z"
const extractCmd = "x"
const extractWoFolders = "e"
const archiveCmd = "a"
const destFolderTemplate = "-o{folder}"
const templatePattern = "{folder}"

/*
-mx：设置 LZMA/LZMA2 压缩算法的字典大小。字典大小是压缩算法用来存储已压缩数据的字典的大小。字典越大，压缩率越高，但压缩速度越慢。
-mfb：设置 LZMA/LZMA2 压缩算法的快速字节查找器 (FB) 的大小。FB 是压缩算法用来查找重复数据的工具。FB 的大小越大，压缩率越高，但压缩速度越慢。
-ms：启用固实模式。固实模式是一种压缩算法，可以提高压缩率，但会降低压缩速度。
-md：设置 LZMA/LZMA2 压缩算法的字典匹配距离。字典匹配距离是压缩算法用来匹配字典中数据的距离。字典匹配距离越大，压缩率越高，但压缩速度越慢。
-myx：设置 LZMA/LZMA2 压缩算法的压缩线程数。压缩线程数是压缩算法用来并行压缩数据的线程数。压缩线程数越多，压缩速度越快，但压缩率可能会降低。
-mtm：禁用多线程压缩。
-mmt：启用多线程压缩。
-mmtf：启用多线程压缩，并使用任务分解。
-md：设置 LZMA/LZMA2 压缩算法的字典大小。
-mmf：设置 LZMA/LZMA2 压缩算法的模式。
-mmc：设置 LZMA/LZMA2 压缩算法的内存限制。
-mpb：禁用进度指示器。
-mlc：禁用日志记录。
*/
const archiveType = "-t7z"
const volumesSize = "-v{size}"
const volumesSizePattern = "{size}"
const archivePassword = "-p{password}"
const passwordPattern = "{password}"

type Options struct {
	Args           []string
	Password       string
	VolumesSize    string
	DestFolderPath string
	Param          string
}

func (o Options) Arg() []string {
	if o.Password != "" {
		o.Args = append(o.Args, strings.Replace(archivePassword, passwordPattern, o.Password, -1))
	}
	if o.VolumesSize != "" {
		o.Args = append(o.Args, strings.Replace(volumesSize, volumesSizePattern, o.VolumesSize, -1))
	}
	if o.DestFolderPath != "" {
		o.Args = append(o.Args, strings.Replace(destFolderTemplate, templatePattern, o.DestFolderPath, -1))
	}
	if o.Param != "" {
		o.Args = append(o.Args, strings.Split(o.Param, "")...)
	}
	return o.Args
}

var CompressOpt = Options{
	Args: []string{
		"-mx=5",
	},
}
