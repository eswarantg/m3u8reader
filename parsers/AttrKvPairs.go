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

type AttrKVPairs struct {
	m map[common.AttrId]interface{}
}

func (a *AttrKVPairs) Done() {
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
		ret += fmt.Sprintf("[ %v[%v]:\"%v\" ],", common.AttrNames[k], k, v)
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
	var ok bool
	var vFloat float64
	var vInt int64
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v not found", common.AttrNames[k])
		return
	}
	vFloat, ok = val.(float64)
	if ok {
		ret = vFloat
	}
	vInt, ok = val.(int64)
	if ok {
		ret = float64(vInt)
	}
	err = fmt.Errorf("%v expected float64 found of data type %v.%v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetInt64(t common.TagId, k common.AttrId) (ret int64, err error) {
	var ok bool
	var vInt int64
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v not found", common.AttrNames[k])
		return
	}
	vInt, ok = val.(int64)
	if ok {
		ret = vInt
	}
	err = fmt.Errorf("%v expected int64 found of data type %v:%v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetTime(t common.TagId, k common.AttrId) (ret time.Time, err error) {
	var ok bool
	var vTime time.Time
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v not found", common.AttrNames[k])
		return
	}
	vTime, ok = val.(time.Time)
	if ok {
		ret = vTime
	}
	err = fmt.Errorf("%v expected time.Time found of data type %v:%v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}

func (a *AttrKVPairs) GetString(t common.TagId, k common.AttrId) (ret string, err error) {
	var ok bool
	var vStr string
	val := a.Get(k)
	if val == nil {
		err = fmt.Errorf("%v not found", common.AttrNames[k])
		return
	}
	vStr, ok = val.(string)
	if ok {
		ret = vStr
	}
	err = fmt.Errorf("%v expected time.Time found of data type %v:%v", common.TagNames[t], common.AttrNames[k], reflect.ValueOf(val).Kind())
	return
}
