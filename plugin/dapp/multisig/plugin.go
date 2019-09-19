package multisig

import (
	"github.com/33cn/chain33/pluginmgr"
	_ "github.com/GM-Publicchain/gm/plugin/dapp/multisig/autotest" //register auto test
	"github.com/GM-Publicchain/gm/plugin/dapp/multisig/commands"
	"github.com/GM-Publicchain/gm/plugin/dapp/multisig/executor"
	"github.com/GM-Publicchain/gm/plugin/dapp/multisig/rpc"
	mty "github.com/GM-Publicchain/gm/plugin/dapp/multisig/types"
	_ "github.com/GM-Publicchain/gm/plugin/dapp/multisig/wallet" // register wallet package
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     mty.MultiSigX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.MultiSigCmd,
		RPC:      rpc.Init,
	})
}
