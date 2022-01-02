package parsers

import (
	"io"

	"github.com/eswarantg/m3u8reader/common"
)

type AttrKVPairs map[common.AttrId]interface{}

func (a *AttrKVPairs) Init() {
	*a = make(map[common.AttrId]interface{})
}
func (a *AttrKVPairs) Store(k common.AttrId, v interface{}) {
	if a == nil {
		a.Init()
	}
	(*a)[k] = v
}
func (a *AttrKVPairs) Clear() {
	for k := range *a {
		delete(*a, k)
	}
}

type M3u8Handler interface {
	PostRecord(tag common.TagId, kvpairs AttrKVPairs) error
}

type Parser interface {
	ParseData(data []byte, handler M3u8Handler) (nBytes int, err error)
	Parse(rdr io.Reader, handler M3u8Handler) (nBytes int, err error)
}
