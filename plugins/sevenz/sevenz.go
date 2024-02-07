package sevenz

import (
	"fmt"
	"os/exec"
	"strings"
)

type Sevenz struct {
}

func (sz Sevenz) Compress(srcPath []string, destPath string) error {
	return CompressInsane(destPath, srcPath, CompressOpt)
}
func (sz Sevenz) CompressWithPwd(srcPath []string, destPath, password string) error {
	opt := CompressOpt
	opt.Password = password
	return CompressInsane(destPath, srcPath, opt)
}
func (sz Sevenz) CompressVolumes(srcPath []string, destPath string, size string) error {
	opt := CompressOpt
	opt.VolumesSize = size
	return CompressInsane(destPath, srcPath, opt)
}
func (sz Sevenz) CompressWithPwdVolumes(srcPath []string, destPath, password string, size string) error {
	opt := CompressOpt
	opt.Password = password
	opt.VolumesSize = size
	return CompressInsane(destPath, srcPath, opt)
}

func (sz Sevenz) Extract(srcPath, destPath, password string) error {
	return Extract(srcPath, destPath, Options{Password: password})
}

func (sz Sevenz) ExtractAble(srcPath string) bool {
	return strings.HasSuffix(srcPath, ".001") || strings.HasSuffix(srcPath, ".7z")
}

func (sz Sevenz) FileName(srcPath string) string {
	tmp := strings.Split(srcPath, ".")
	if strings.HasSuffix(srcPath, ".001") {
		return strings.Join(tmp[:len(tmp)-2], ".")
	}
	return strings.Join(tmp[:len(tmp)-1], ".")

}

func Extract(archiveFullPath, destFullPath string, opts ...Options) error {
	opt := checkOption(opts)
	opt.DestFolderPath = destFullPath
	opt.Args = append(opt.Args, archiveFullPath)

	cmd := exec.Command(sevenCmd, append([]string{extractCmd}, opt.Arg()...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)
	}
	return err
}

func ExtractWoFolders(archiveFullPath, destFullPath string, opts ...Options) error {
	opt := checkOption(opts)
	opt.DestFolderPath = destFullPath
	opt.Args = append(opt.Args, archiveFullPath)

	cmd := exec.Command(sevenCmd, append([]string{extractWoFolders}, opt.Arg()...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)
	}
	return err
}

func CompressInsane(archiveFullPath string, srcFullPath []string, opts ...Options) error {
	opt := checkOption(opts)
	opt.Args = append(opt.Args, archiveFullPath)
	opt.Args = append(opt.Args, srcFullPath...)
	//cmd := exec.Command(sevenCmd, archiveCmd, archiveType, insaneCompressionParams, archiveFullPath, srcFullPath)
	cmd := exec.Command(sevenCmd, append([]string{archiveCmd, archiveType}, opt.Arg()...)...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)
	}
	return err
}

func checkOption(opts []Options) Options {
	if len(opts) > 0 {
		return opts[0]
	}
	return Options{}
}
