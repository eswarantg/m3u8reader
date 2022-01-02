package m3u8reader

import (
	"fmt"

	"github.com/eswarantg/m3u8reader/common"
)

type M3U8Entry struct {
	Tag    common.TagId
	Values map[common.AttrId]interface{}
}

func (m *M3U8Entry) StoreKV(k common.AttrId, v interface{}) {
	if m.Values == nil {
		m.Values = make(map[common.AttrId]interface{})
	}
	m.Values[k] = v
}

func (m *M3U8Entry) String() string {
	return fmt.Sprintf("%v %v", m.Tag, m.Values)
}

func (m *M3U8Entry) URI() (string, error) {
	switch m.Tag {
	case common.M3U8ExtXStreamInf:
		return m.Values[common.INTUnknownAttr].(string), nil
	case common.M3U8ExtXMedia:
		return m.Values[common.M3U8Uri].(string), nil
	case common.M3U8ExtInf:
		return m.Values[common.M3U8Uri].(string), nil
	case common.M3U8ExtXPreLoadHint:
		return m.Values[common.M3U8Uri].(string), nil
	}
	return "", fmt.Errorf("URI not available")
}
