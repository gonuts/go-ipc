package common

import (
	"bufio"
	"fmt"
	"io"
)

type PluginIO struct {
	io.Reader
	io.Writer
	io.Closer
}

func NewPluginIO(stdin io.ReadCloser, stdout io.WriteCloser) PluginIO {

	return PluginIO{stdin, stdout, stdin}
}

func StdErrForward(stderr io.Reader) {
	bio := bufio.NewReader(stderr)

	for {
		line, _, err := bio.ReadLine()
		if err == nil {
			fmt.Println("plugin:", string(line))
		} else {
			fmt.Println(err)
			break
		}
	}
}
