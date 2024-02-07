package plugins

type IPlugin interface {
	ICompresser
	IExtracter
}

type ICompresser interface {
	Compress(srcPath []string, destPath string) error
	CompressWithPwd(srcPath []string, destPath, password string) error
	CompressVolumes(srcPath []string, destPath string, size string) error
	CompressWithPwdVolumes(srcPath []string, destPath, password string, size string) error
}

type IExtracter interface {
	Extract(srcPath, destPath, password string) error
	ExtractAble(srcPath string) bool
	FileName(srcPath string) string
}
