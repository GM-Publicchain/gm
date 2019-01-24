// +build go1.8

package main

import (
	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin"
	_ "github.com/GM-Publicchain/gm/plugin"

	"github.com/33cn/chain33/util/cli"
)

func main() {
	cli.RunChain33("gm")
}
