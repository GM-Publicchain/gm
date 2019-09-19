// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token 创建token
package token

import (
	"github.com/33cn/chain33/pluginmgr"
	_ "github.com/GM-Publicchain/gm/plugin/dapp/token/autotest" // register token autotest package
	"github.com/GM-Publicchain/gm/plugin/dapp/token/commands"
	"github.com/GM-Publicchain/gm/plugin/dapp/token/executor"
	"github.com/GM-Publicchain/gm/plugin/dapp/token/rpc"
	"github.com/GM-Publicchain/gm/plugin/dapp/token/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.TokenX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.TokenCmd,
		RPC:      rpc.Init,
	})
}
