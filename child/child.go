package child

import (
	"go-plugin/common"
	"go-plugin/common/api"
	"go-plugin/ipc"
	"io"
	"net/rpc"
	"os"
)

type Plugin struct {
	Name string
	RPC  *rpc.Server
	mpx  *ipc.Multiplex
}

func NewPlugin(name string) *Plugin {
	plugin := new(Plugin)
	plugin.Name = name
	plugin.RPC = rpc.NewServer()
	plugin.init()
	pi := new(PluginInfo)
	pi.Plugin = plugin

	plugin.RPC.Register(pi)

	return plugin
}

func (p *Plugin) init() {
	// redirect stdout
	stdout := os.Stdout
	os.Stdout = os.Stderr
	pio := common.NewPluginIO(os.Stdin, stdout)
	p.mpx = ipc.NewMultiplex(pio)
}

func (p *Plugin) Start() {

	rpcreader := p.mpx.RawReaderChannel("rpcw")
	rpcwriter := p.mpx.RawWriterChannel("rpcr")

	go p.RPC.ServeConn(common.NewPluginIO(rpcreader, rpcwriter))
}

type PluginInfo struct {
	Plugin *Plugin
}

func (pi *PluginInfo) GetName(args *api.EmptyArgs, result *string) error {
	*result = pi.Plugin.Name
	return nil
}

func (pc *Plugin) ReaderChannel(name string) io.Reader {
	var tmp interface{}
	tmp = pc.mpx.RawReaderChannel(name)
	return tmp.(io.Reader)
}
func (pc *Plugin) WriterChannel(name string) io.Writer {
	var tmp interface{}
	tmp = pc.mpx.RawWriterChannel(name)
	return tmp.(io.Writer)
}
