// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evm

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/GM-Publicchain/gm/plugin/dapp/evm/commands"
	"github.com/GM-Publicchain/gm/plugin/dapp/evm/executor"
	"github.com/GM-Publicchain/gm/plugin/dapp/evm/rpc"
	"github.com/GM-Publicchain/gm/plugin/dapp/evm/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.ExecutorName,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.EvmCmd,
		RPC:      rpc.Init,
	})
}
