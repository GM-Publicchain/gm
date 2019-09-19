// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package privacy

import (
	"github.com/33cn/chain33/pluginmgr"
	_ "github.com/GM-Publicchain/gm/plugin/dapp/privacy/autotest" // register autotest package
	"github.com/GM-Publicchain/gm/plugin/dapp/privacy/commands"
	_ "github.com/GM-Publicchain/gm/plugin/dapp/privacy/crypto" // register crypto package
	"github.com/GM-Publicchain/gm/plugin/dapp/privacy/executor"
	"github.com/GM-Publicchain/gm/plugin/dapp/privacy/rpc"
	"github.com/GM-Publicchain/gm/plugin/dapp/privacy/types"
	_ "github.com/GM-Publicchain/gm/plugin/dapp/privacy/wallet" // register wallet package
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.PrivacyX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.PrivacyCmd,
		RPC:      rpc.Init,
	})
}
