package ipc

import (
	"fmt"
	"os"
	"syscall"
)

type SocketPair struct {
	fd         []int
	LocalFile  *os.File
	RemoteFile *os.File
}

func NewSocketPair() *SocketPair {
	sp := new(SocketPair)
	sp.init()
	return sp
}

func (sp *SocketPair) init() {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)

	if err != nil {
		fmt.Println(err)
		return
	}
	sp.LocalFile = os.NewFile(uintptr(fds[0]), "sp-0")
	sp.RemoteFile = os.NewFile(uintptr(fds[1]), "sp-1")

}
func (sp *SocketPair) Test() {

}
