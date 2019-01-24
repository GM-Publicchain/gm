package plugin

import (
	_ "github.com/GM-Publicchain/gm/plugin/consensus/init" //consensus init
	_ "github.com/GM-Publicchain/gm/plugin/crypto/init"    //crypto init
	_ "github.com/GM-Publicchain/gm/plugin/dapp/init"      //dapp init
	_ "github.com/GM-Publicchain/gm/plugin/mempool/init"   //mempool init
	_ "github.com/GM-Publicchain/gm/plugin/store/init"     //store init
)
