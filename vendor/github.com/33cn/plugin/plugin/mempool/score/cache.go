package score

import (
	"bytes"
	"encoding/gob"

	"github.com/33cn/chain33/common/skiplist"
	"github.com/33cn/chain33/system/mempool"
	"github.com/33cn/chain33/types"
)

var mempoolDupResendInterval int64 = 600 // mempool内交易过期时间，10分钟

// Queue 分数队列模式(分数=常量a*手续费/交易字节数-常量b*时间*定量c,按分数排队,高的优先,常量a,b和定量c可配置)
type Queue struct {
	txMap     map[string]*skiplist.SkipValue
	txList    *skiplist.SkipList
	subConfig subConfig
}

// NewQueue 创建队列
func NewQueue(subcfg subConfig) *Queue {
	return &Queue{
		txMap:     make(map[string]*skiplist.SkipValue, subcfg.PoolCacheSize),
		txList:    skiplist.NewSkipList(&skiplist.SkipValue{Score: -1, Value: nil}),
		subConfig: subcfg,
	}
}

func (cache *Queue) newSkipValue(item *mempool.Item) (*skiplist.SkipValue, error) {
	//tx := item.value
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(item.Value)
	if err != nil {
		return nil, err
	}
	size := len(buf.Bytes())
	return &skiplist.SkipValue{Score: cache.subConfig.PriceConstant*(item.Value.Fee/int64(size))*cache.subConfig.PricePower - cache.subConfig.TimeParam*item.EnterTime, Value: item}, nil
}

// Exist 是否存在
func (cache *Queue) Exist(hash string) bool {
	_, exists := cache.txMap[hash]
	return exists
}

//GetItem 获取数据通过 key
func (cache *Queue) GetItem(hash string) (*mempool.Item, error) {
	if k, exist := cache.txMap[hash]; exist {
		return k.Value.(*mempool.Item), nil
	}
	return nil, types.ErrNotFound
}

// Push 把给定tx添加到Queue；如果tx已经存在Queue中或Mempool已满则返回对应error
func (cache *Queue) Push(item *mempool.Item) error {
	hash := item.Value.Hash()
	if cache.Exist(string(hash)) {
		s := cache.txMap[string(hash)]
		addedItem := s.Value.(*mempool.Item)
		addedTime := addedItem.EnterTime

		if types.Now().Unix()-addedTime < mempoolDupResendInterval {
			return types.ErrTxExist
		}
		// 超过2分钟之后的重发交易返回nil，再次发送给P2P，但是不再次加入mempool
		// 并修改其enterTime，以避免该交易一直在节点间被重发
		newEnterTime := types.Now().Unix()
		resendItem := &mempool.Item{Value: item.Value, Priority: item.Value.Fee, EnterTime: newEnterTime}
		var err error
		sv, err := cache.newSkipValue(resendItem)
		if err != nil {
			return err
		}
		cache.Remove(string(hash))
		cache.txList.Insert(sv)
		cache.txMap[string(hash)] = sv
		// ------------------
		return nil
	}

	it := &mempool.Item{Value: item.Value, Priority: item.Value.Fee, EnterTime: item.EnterTime}
	sv, err := cache.newSkipValue(it)
	if err != nil {
		return err
	}
	if int64(cache.txList.Len()) >= cache.subConfig.PoolCacheSize {
		tail := cache.txList.GetIterator().Last()
		//分数高存留
		if sv.Compare(tail) == -1 {
			cache.Remove(string(tail.Value.(*mempool.Item).Value.Hash()))
		} else {
			return types.ErrMemFull
		}
	}
	cache.txList.Insert(sv)
	cache.txMap[string(hash)] = sv
	return nil
}

// Remove 删除数据
func (cache *Queue) Remove(hash string) error {
	cache.txList.Delete(cache.txMap[hash])
	delete(cache.txMap, hash)
	return nil
}

// Size 数据总数
func (cache *Queue) Size() int {
	return cache.txList.Len()
}

// Walk 遍历整个队列
func (cache *Queue) Walk(count int, cb func(value *mempool.Item) bool) {
	i := 0
	cache.txList.Walk(func(item interface{}) bool {
		if !cb(item.(*mempool.Item)) {
			return false
		}
		i++
		return i != count
	})
}
