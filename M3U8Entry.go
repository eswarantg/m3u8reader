package m3u8reader

import (
	"fmt"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type M3U8Entry struct {
	Tag    common.TagId
	Values *parsers.AttrKVPairs
}

func (m *M3U8Entry) Done() {
	m.Values.Done()
	m.Values = nil
}

func (m *M3U8Entry) StoreKV(k common.AttrId, v interface{}) {
	m.Values.Store(k, v)
}

func (m *M3U8Entry) String() string {
	return fmt.Sprintf("%v(%v)%v", common.TagNames[m.Tag], m.Tag, m.Values.String())
}

func (m *M3U8Entry) URI() (string, error) {
	var attrId common.AttrId = -1
	switch m.Tag {
	case common.M3U8ExtXStreamInf:
		attrId = common.INTUnknownAttr
	case common.M3U8ExtXMedia, common.M3U8ExtInf, common.M3U8ExtXPreLoadHint, common.M3U8ExtXPart:
		attrId = common.M3U8Uri
	}
	if attrId != -1 {
		val := m.Values.Get(attrId)
		if val != nil {
			switch v := val.(type) {
			case string:
				return v, nil
			case []byte:
				return string(v), nil
			}
		}
	}
	return "", fmt.Errorf("URI not available")
}
