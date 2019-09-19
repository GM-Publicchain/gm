package js

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/GM-Publicchain/gm/plugin/dapp/js/executor"
	ptypes "github.com/GM-Publicchain/gm/plugin/dapp/js/types"

	// init auto test
	_ "github.com/GM-Publicchain/gm/plugin/dapp/js/autotest"
	"github.com/GM-Publicchain/gm/plugin/dapp/js/command"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     ptypes.JsX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      command.JavaScriptCmd,
		RPC:      nil,
	})
}
