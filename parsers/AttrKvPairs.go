package parsers

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/eswarantg/m3u8reader/common"
)

var AttrKVPairsPool = sync.Pool{
	New: func() interface{} {
		ret := &AttrKVPairs{
			m: make(map[common.AttrId]interface{}, 5),
		}
		return ret
	},
}

//Initialized to ZERO/FALSE....automatically
var AttrKVPairsSyncPool bool

type AttrKVPairs struct {
	m map[common.AttrId]interface{}
}

func (a *AttrKVPairs) Done() {
	if !AttrKVPairsSyncPool {
		return
	}
	for k := range a.m {
		delete(a.m, k)
	}
	if len(a.m) != 0 {
		panic("\nPutting back into AttrKVPairsPool without full cleanup")
	}
	AttrKVPairsPool.Put(a)
	a = nil
}

func NewAttrKVPairsDebug(label string) (ret *AttrKVPairs) {
	if !AttrKVPairsSyncPool {
		return &AttrKVPairs{m: make(map[common.AttrId]interface{}, 10)}
	}
	var ok bool
	obj := AttrKVPairsPool.Get()
	if obj == nil {
		panic(fmt.Sprintf("\nAttrKVPairsPool returned nil : invoked at %v", label))
		//return
	}
	ret, ok = obj.(*AttrKVPairs)
	if !ok {
		panic(fmt.Sprintf("\nAttrKVPairsPool returned bad object type %v invoked at %v", reflect.TypeOf(obj).Kind(), label))
		//return
	}
	//fmt.Printf("NewAttrKVPairsDebug return valid for %v\n", label)
	return ret
}

func NewAttrKVPairs() *AttrKVPairs {
	return NewAttrKVPairsDebug("NewAttrKVPairs")
}

func (a *AttrKVPairs) Map() map[common.AttrId]interface{} {
	if a == nil {
		return nil
	}
	return a.m
}

func (a *AttrKVPairs) Get(k common.AttrId) (v interface{}) {
	var ok bool
	if a == nil {
		return
	}
	v, ok = a.m[k]
	if !ok {
		v = nil
	}
	return
}
func (a *AttrKVPairs) String() string {
	if a == nil {
		return "<nil"
	}
	ret := ""
	for k, v := range (*a).m {
		switch val := v.(type) {
		case []byte:
			ret += fmt.Sprintf("[ %v[%v]:\"%v\" ],", common.AttrNames[k], k, string(val))
		default:
			ret += fmt.Sprintf("[ %v[%v]:\"%v\" ],", common.AttrNames[k], k, val)
		}
	}
	return ret
}

func (a *AttrKVPairs) Exists(k common.AttrId) bool {
	var ok bool
	if a == nil {
		return false
	}
	_, ok = a.m[k]
	return ok
}
func (a *AttrKVPairs) StoreDebug(label string, k common.AttrId, v interface{}) {
	if a == nil {
		panic(fmt.Sprintf("\nAttrKVPairs not allocated at %v", label))
	}
	a.m[k] = v
}

func (a *AttrKVPairs) Store(k common.AttrId, v interface{}) {
	a.StoreDebug("AttrKVPairs.Store", k, v)
}

func (a *AttrKVPairs) GetFloat64(t common.TagId, k common.AttrId) (ret float64, err error) {
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v:%v not found", common.TagNames[t], common.AttrNames[k])
		return
	}
	switch v := val.(type) {
	case float64:
		return v, nil
	case int64:
		return float64(v), nil
	}
	err = fmt.Errorf("%v:%v expected float64 found of data type %v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetInt64(t common.TagId, k common.AttrId) (ret int64, err error) {
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v:%v not found", common.TagNames[t], common.AttrNames[k])
		return
	}
	switch v := val.(type) {
	case int64:
		return v, nil
	}
	err = fmt.Errorf("%v:%v expected int64 found of data type %v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetTime(t common.TagId, k common.AttrId) (ret time.Time, err error) {
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v:%v not found", common.TagNames[t], common.AttrNames[k])
		return
	}
	switch v := val.(type) {
	case time.Time:
		return v, nil
	}
	err = fmt.Errorf("%v:%v expected time.Time found of data type %v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetByteRange(t common.TagId, k common.AttrId) (ret [2]int64, err error) {
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v:%v not found", common.TagNames[t], common.AttrNames[k])
		return
	}
	switch v := val.(type) {
	case [2]int64:
		return v, nil
	}
	err = fmt.Errorf("%v:%v expected [2]int64 found of data type %v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetString(t common.TagId, k common.AttrId) (ret string, err error) {
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v:%v not found", common.TagNames[t], common.AttrNames[k])
		return
	}
	switch v := val.(type) {
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	}
	err = fmt.Errorf("%v:%v expected time.Time found of data type %v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}
