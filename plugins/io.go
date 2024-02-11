package plugins

import (
	"bufio"
	"fmt"
	"io"
)

func HandleStdout(out io.Reader) {
	in := bufio.NewScanner(out)
	for in.Scan() {
		fmt.Println(in.Text())
	}
}
