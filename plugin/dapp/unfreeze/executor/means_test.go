package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common"
	pty "github.com/GM-Publicchain/gm/plugin/dapp/unfreeze/types"
	"fmt"
	"time"
	"google.golang.org/grpc"
	"context"
	"bytes"
)

func TestCalcFrozen(t *testing.T) {
	types.SetTitleOnlyForTest("chain33")
	m, err := newMeans("LeftProportion", 15000000)
	assert.Nil(t, err)
	assert.NotNil(t, m)

	cases := []struct {
		start         int64
		now           int64
		period        int64
		total         int64
		tenThousandth int64
		expect        int64
	}{
		{10000, 10001, 10, 10000, 2, 9998},
		{10000, 10011, 10, 10000, 2, 9996},
		{10000, 10001, 10, 1e17, 2, 9998 * 1e13},
		{10000, 10011, 10, 1e17, 2, 9998 * 9998 * 1e9},
	}

	for _, c := range cases {
		c := c
		t.Run("test LeftProportion", func(t *testing.T) {
			create := pty.UnfreezeCreate{
				StartTime:   c.start,
				AssetExec:   "coins",
				AssetSymbol: "bty",
				TotalCount:  c.total,
				Beneficiary: "x",
				Means:       "LeftProportion",
				MeansOpt: &pty.UnfreezeCreate_LeftProportion{
					LeftProportion: &pty.LeftProportion{
						Period:        c.period,
						TenThousandth: c.tenThousandth,
					},
				},
			}
			u := &pty.Unfreeze{
				TotalCount: c.total,
				Means:      "LeftProportion",
				StartTime:  c.start,
				MeansOpt: &pty.Unfreeze_LeftProportion{
					LeftProportion: &pty.LeftProportion{
						Period:        c.period,
						TenThousandth: c.tenThousandth,
					},
				},
			}
			u, err = m.setOpt(u, &create)
			assert.Nil(t, err)

			f, err := m.calcFrozen(u, c.now)
			assert.Nil(t, err)

			assert.Equal(t, c.expect, f)

		})
	}
}

func TestLeftV1(t *testing.T) {
	cases := []struct {
		start         int64
		now           int64
		period        int64
		total         int64
		tenThousandth int64
		expect        int64
	}{
		{10000, 10001, 10, 10000, 2, 9998},
		{10000, 10011, 10, 10000, 2, 9996},
		{10000, 10001, 10, 1e17, 2, 9998 * 1e13},
		{10000, 10011, 10, 1e17, 2, 9998 * 9998 * 1e9},
	}

	for _, c := range cases {
		c := c
		t.Run("test LeftProportionV1", func(t *testing.T) {
			create := pty.UnfreezeCreate{
				StartTime:   c.start,
				AssetExec:   "coins",
				AssetSymbol: "bty",
				TotalCount:  c.total,
				Beneficiary: "x",
				Means:       pty.LeftProportionX,
				MeansOpt: &pty.UnfreezeCreate_LeftProportion{
					LeftProportion: &pty.LeftProportion{
						Period:        c.period,
						TenThousandth: c.tenThousandth,
					},
				},
			}
			u := &pty.Unfreeze{
				TotalCount: c.total,
				Means:      pty.LeftProportionX,
				StartTime:  c.start,
				MeansOpt: &pty.Unfreeze_LeftProportion{
					LeftProportion: &pty.LeftProportion{
						Period:        c.period,
						TenThousandth: c.tenThousandth,
					},
				},
			}
			m := leftProportion{}
			u, err := m.setOpt(u, &create)
			assert.Nil(t, err)

			f, err := m.calcFrozen(u, c.now)
			assert.Nil(t, err)

			assert.Equal(t, c.expect, f)

		})
	}
}

func TestFixV1(t *testing.T) {
	cases := []struct {
		start  int64
		now    int64
		period int64
		total  int64
		amount int64
		expect int64
	}{
		{10000, 10001, 10, 10000, 2, 9998},
		{10000, 10011, 10, 10000, 2, 9996},
		{10000, 10001, 10, 1e17, 2, 1e17 - 2},
		{10000, 10011, 10, 1e17, 2, 1e17 - 4},
	}

	for _, c := range cases {
		c := c
		t.Run("test FixAmountV1", func(t *testing.T) {
			create := pty.UnfreezeCreate{
				StartTime:   c.start,
				AssetExec:   "coins",
				AssetSymbol: "bty",
				TotalCount:  c.total,
				Beneficiary: "x",
				Means:       pty.FixAmountX,
				MeansOpt: &pty.UnfreezeCreate_FixAmount{
					FixAmount: &pty.FixAmount{
						Period: c.period,
						Amount: c.amount,
					},
				},
			}
			u := &pty.Unfreeze{
				TotalCount: c.total,
				Means:      pty.FixAmountX,
				StartTime:  c.start,
				MeansOpt: &pty.Unfreeze_FixAmount{
					FixAmount: &pty.FixAmount{
						Period: c.period,
						Amount: c.amount,
					},
				},
			}
			m := fixAmount{}
			u, err := m.setOpt(u, &create)
			assert.Nil(t, err)

			f, err := m.calcFrozen(u, c.now)
			assert.Nil(t, err)

			assert.Equal(t, c.expect, f)

		})
	}
}

// 查询可提币量， 和当前时间有关， 如对应节点时间不对， 查询结果也不对
func TestLeftV2(t *testing.T) {
	cases := []struct {
		start         int64
		now           int64
		period        int64
		total         int64
		tenThousandth int64
		expect        int64
	}{
		{1561607389, 1561607389 + 500000, 67200, 11111130, 1, 11102244},
		{1561607389, -156107389 + 500000, 67200, 11111130, 1, 11111130},
	}

	for _, c := range cases {
		c := c
		t.Run("test LeftProportionV2", func(t *testing.T) {
			create := pty.UnfreezeCreate{
				StartTime:   c.start,
				AssetExec:   "coins",
				AssetSymbol: "bty",
				TotalCount:  c.total,
				Beneficiary: "x",
				Means:       pty.LeftProportionX,
				MeansOpt: &pty.UnfreezeCreate_LeftProportion{
					LeftProportion: &pty.LeftProportion{
						Period:        c.period,
						TenThousandth: c.tenThousandth,
					},
				},
			}
			u := &pty.Unfreeze{
				TotalCount: c.total,
				Means:      pty.LeftProportionX,
				StartTime:  c.start,
				MeansOpt: &pty.Unfreeze_LeftProportion{
					LeftProportion: &pty.LeftProportion{
						Period:        c.period,
						TenThousandth: c.tenThousandth,
					},
				},
			}
			m := leftProportion{}
			u, err := m.setOpt(u, &create)
			assert.Nil(t, err)

			f, err := m.calcFrozen(u, c.now)
			assert.Nil(t, err)

			assert.Equal(t, c.expect, f)

		})
	}
}

func TestDecreaseAmount(t *testing.T) {
	amount := int64(400000)
	rate := int64(1000)

	start := int64(1546272000)
	total := 800000000
	period := 86400 //测试 decreasePeriod设置 period的整数倍
	decreasePeriod := 200*86400
	firstAmount := 400000
	decreaseNum := 34

	t.Run("test TestDecreaseAmount", func(t *testing.T) {
		create := pty.UnfreezeCreate{
			StartTime:   start,
			AssetExec:   "coins",
			AssetSymbol: "bty",
			TotalCount:  int64(total),
			Beneficiary: "x",
			Means:       pty.DecreaseAmountX,
			MeansOpt: &pty.UnfreezeCreate_DecreaseAmount{
				DecreaseAmount: &pty.DecreaseAmount{
					Period:        int64(period),
					TenThousandth: rate,
					FirstDecreaseAmount:int64(firstAmount),
					DecreasePeriod:int64(decreasePeriod),
					DecreaseNums:int64(decreaseNum),
				},
			},
		}
		u := &pty.Unfreeze{
			TotalCount: int64(total),
			Means:      pty.DecreaseAmountX,
			StartTime:  start,
			MeansOpt: &pty.Unfreeze_DecreaseAmount{
				DecreaseAmount: &pty.DecreaseAmount{
					Period:        int64(period),
					TenThousandth: rate,
					FirstDecreaseAmount:int64(firstAmount),
					DecreasePeriod:int64(decreasePeriod),
					DecreaseNums:int64(decreaseNum),
				},
			},
		}
		m := decreaseAmount{}
		u, err := m.setOpt(u, &create)
		if err != nil {
			panic(err)
		}
		days := int64(0)
		//assert.Nil(t, err)
		for i:= int64(1);i <= int64(decreaseNum+11); i ++ {
			tt := getDecreasePeriodAmount(int64(i-1),amount,rate)

			fmt.Println("amount:",tt)
			for n:= int64(1);n <= int64(decreasePeriod/period);n++ {
				days ++
				f, err := m.calcFrozen(u, start + (days-1)*int64(period))
				//assert.Nil(t, err)
				if err != nil {
					panic(err)
				}
				fmt.Println("decrease nums:",i-1,"left amount:",f,time.Unix((start + (days-1)*int64(period)),0).Format("2006-01-02 15:04:05"))
			}
		}

		//assert.Equal(t, c.expect, f)

	})
}


func TestCompareV1V2(t *testing.T) {

	v1grpc := "47.90.48.67:8902"
	v2grpc := "47.90.48.67:8802"
	var v1Client,v2Client types.Chain33Client
	maxReceLimit := grpc.WithMaxMsgSize(30*1024*1024)
	conn, err := grpc.Dial(v1grpc, grpc.WithInsecure(),maxReceLimit)
	if err != nil {
		panic(err)
	}
	v1Client = types.NewChain33Client(conn)

	conn, err = grpc.Dial(v2grpc, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	v2Client = types.NewChain33Client(conn)
	maxHeight ,err := v1Client.GetLastHeader(context.Background(),&types.ReqNil{})
	if err != nil {
		panic(err)
	}
	fmt.Println("MaxHeight",maxHeight)
	//version1,_ := v1Client.Version(context.Background(),&types.ReqNil{})
	//fmt.Println("version1",version1)
	//version2,_ := v2Client.Version(context.Background(),&types.ReqNil{})
	//fmt.Println("version2",version2)
	//txs1,_ := v1Client.GetTransactionByAddr(context.Background(),&types.ReqAddr{Addr:"1P7P4v3kL39zugQgDDLRqxzGjQd7aEbfKs",Direction:0,Count:250,Height:-1,Index:0})
	////fmt.Println("txs1",txs1)
	//txs2,_ := v1Client.GetTransactionByAddr(context.Background(),&types.ReqAddr{Addr:"1P7P4v3kL39zugQgDDLRqxzGjQd7aEbfKs",Direction:0,Count:250,Height:-1,Index:0})
	////fmt.Println("txs2",txs2)
	//for n:=1 ;n<= len(txs1.TxInfos) ;n++{
	//	h1 := txs1.TxInfos[n].Hash
	//	h2 := txs2.TxInfos[n].Hash
	//
	//	r1,_:= v1Client.QueryTransaction(context.Background(),&types.ReqHash{Hash:h1})
	//	fmt.Println("h1 ty ",r1.Receipt.Ty)
	//	r2,_:= v2Client.QueryTransaction(context.Background(),&types.ReqHash{Hash:h2})
	//	fmt.Println("h2 ty ",r2.Receipt.Ty)
	//	if r1.Receipt.Ty != r2.Receipt.Ty {
	//		fmt.Println("h1 hash",common.ToHex(h1))
	//		fmt.Println("h2 hash",common.ToHex(h2))
	//	}
	//}

	go func() {
		for n := int64(5000);n <= maxHeight.Height;n++ {
			fmt.Println("height:",n)
			v1,err:= v1Client.GetBlocks(context.Background(),&types.ReqBlocks{Start:n,End:n,IsDetail:true})
			if err != nil {
				fmt.Println("v1 err",err)
				break
			}
			if !v1.IsOk {
				fmt.Println("v1 not ok")
				break
			}
			var v1Block,v2Block types.BlockDetails
			err = types.Decode(v1.Msg,&v1Block)
			if err != nil {
				panic(err)
			}
			v2,err:= v2Client.GetBlocks(context.Background(),&types.ReqBlocks{Start:n,End:n,IsDetail:true})
			if err != nil {
				fmt.Println("v2 err",err)
				break
			}
			if !v2.IsOk {
				fmt.Println("v2 not ok")
				break
			}
			err = types.Decode(v2.Msg,&v2Block)
			if err != nil {
				panic(err)
			}
			//fmt.Println("v1",v1Block)
			//fmt.Println("v2",v2Block)
			if !bytes.Equal(v1Block.Items[0].Block.TxHash,v2Block.Items[0].Block.TxHash) {
				fmt.Println("err txhash note equal")
				break
			}
			fmt.Println("v1 txhash",common.ToHex(v1Block.Items[0].Block.TxHash)," v2 txhash",common.ToHex(v2Block.Items[0].Block.TxHash))
			if !bytes.Equal(v1Block.Items[0].Block.StateHash,v2Block.Items[0].Block.StateHash) {
				fmt.Println("err StateHash note equal")
				break
			}
			fmt.Println("v1 StateHash",common.ToHex(v1Block.Items[0].Block.StateHash)," v2 StateHash",common.ToHex(v2Block.Items[0].Block.StateHash))
			if !bytes.Equal(v1Block.Items[0].Block.ParentHash,v2Block.Items[0].Block.ParentHash) {
				fmt.Println("err ParentHash note equal")
				break
			}
			fmt.Println("v1 ParentHash",common.ToHex(v1Block.Items[0].Block.ParentHash)," v2 ParentHash",common.ToHex(v2Block.Items[0].Block.ParentHash))
		}
	}()
	select {
	
	}
}