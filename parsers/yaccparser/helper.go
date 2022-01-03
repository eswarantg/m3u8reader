package yaccparser

import (
	"fmt"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type keyValuePair struct {
	k common.AttrId
	v interface{}
}

type keyValuePairs struct {
	kvs *parsers.AttrKVPairs
}

func (e *keyValuePairs) storeKVDebug(label string, k common.AttrId, v interface{}) {
	//fmt.Printf("%v keyValuePairs.storeKVDebug called for %v %v %v\n", e, label, common.AttrNames[k], v)
	if e.kvs == nil {
		t := parsers.NewAttrKVPairsDebug(label + "_keyValuePairs.storeKVDebug")
		e.kvs = t
	}
	if e.kvs == nil {
		panic(fmt.Sprintf("keyValuePairs.storeKVDebug allocated but still nil %v", label))
	}
	e.kvs.StoreDebug(label, k, v)
	//fmt.Printf("%v keyValuePairs.storeKVDebug done for %v\n", e, label)
}
func (e *keyValuePairs) clear(label string) {
	//fmt.Printf("%v keyValuePairs.clear called for %v\n", e, label)
	e.kvs = nil
	//fmt.Printf("%v keyValuePairs.clear done for %v\n", e, label)
}

type accEntry struct {
	tag common.TagId
	kvs *parsers.AttrKVPairs
}

func (e *accEntry) assignKVPS(label string, val keyValuePairs) {
	//fmt.Printf("%v accEntry.assignKVPS called for %v\n", e, label)
	e.kvs = val.kvs
	val.kvs = nil
	//fmt.Printf("%v accEntry.assignKVPS done for %v\n", e, label)
}
func (e *accEntry) storeKVDebug(label string, k common.AttrId, v interface{}) {
	//fmt.Printf("%v accEntry.storeKVDebug called for %v %v %v\n", e, label, common.AttrNames[k], v)
	if e.kvs == nil {
		t := parsers.NewAttrKVPairsDebug(label + "_accEntry.storeKVDebug")
		e.kvs = t
	}
	if e.kvs == nil {
		panic(fmt.Sprintf("accEntry.storeKVDebug allocated but still nil %v", label))
	}
	e.kvs.StoreDebug(label, k, v)
	//fmt.Printf("%v accEntry.storeKVDebug done for %v\n", e, label)
}
func (e *accEntry) clear(label string) {
	//fmt.Printf("%v accEntry.clear called for %v\n", e, label)
	e.kvs = nil
	//fmt.Printf("%v accEntry.clear done for %v\n", e, label)
}
