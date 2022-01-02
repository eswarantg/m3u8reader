package yaccparser

import (
	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type keyValuePair struct {
	k common.AttrId
	v interface{}
}

type accEntry struct {
	tag common.TagId
	kvs parsers.AttrKVPairs
}

func (e *accEntry) assignKVPS(val parsers.AttrKVPairs) {
	e.kvs = val
}
func (e *accEntry) storeKV(k common.AttrId, v interface{}) {
	if e.kvs == nil {
		e.kvs = parsers.AttrKVPairs{}
		e.kvs.Init()
	}
	e.kvs.Store(k, v)
}
