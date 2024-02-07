package woker

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

var tasks = []Task{}
var bar *progressbar.ProgressBar
var cnt = 0

func RunBar() {
	bar = progressbar.NewOptions(len(tasks),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for _, t := range tasks {
		Add(t)
	}
	Finished()
}

func Process(t Task) {
	cnt++
	bar.Describe(fmt.Sprintf("[cyan][%d/%d][reset] %s", cnt, bar.GetMax(), t.Src))
	bar.Add(1)
}

func AddTask(t Task) {
	tasks = append(tasks, t)
}
