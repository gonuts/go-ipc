package parent

import (
	"fmt"
	"go-plugin/common"
	"go-plugin/ipc"
	"net/rpc"
	"os"
	"os/exec"
)

type PluginClient struct {
	File string
	proc *exec.Cmd
	RPC  *rpc.Client
	mpx  *ipc.Multiplex
}

func LoadPlugin(file string) *PluginClient {
	pc := new(PluginClient)
	pc.File = file
	return pc
}

func (pc *PluginClient) Start() error {
	if pc.proc != nil {
		return common.NewPluginError("plugin is already started")
	}

	pc.proc = exec.Command(pc.File)
	stderr, _ := pc.proc.StderrPipe()
	stdin, _ := pc.proc.StdinPipe()
	stdout, _ := pc.proc.StdoutPipe()

	fmt.Printf("Starting plugin %s\n", pc.File)
	err := pc.proc.Start()

	if err != nil {
		return err
	}

	go common.StdErrForward(stderr)

	pio := common.NewPluginIO(stdout, stdin)
	pc.mpx = ipc.NewMultiplex(pio)

	rpcwriter := pc.mpx.RawWriterChannel("rpcw")
	rpcreader := pc.mpx.RawReaderChannel("rpcr")

	pc.RPC = rpc.NewClient(common.NewPluginIO(rpcreader, rpcwriter))

	if !pc.VerifyRPC() {
		return common.NewPluginError("plugin is not a valid plugin or is not responding")
	}

	return nil
}

func (pc *PluginClient) VerifyRPC() bool {
	return true
}

func (pc *PluginClient) Stop() {
	pc.proc.Process.Signal(os.Interrupt)

}

func (pc *PluginClient) ReaderChannel(name string) *ipc.Channel {
	return pc.mpx.RawReaderChannel(name)
}

func (pc *PluginClient) WriterChannel(name string) *ipc.Channel {
	return pc.mpx.RawWriterChannel(name)
}
