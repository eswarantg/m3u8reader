package parsers

import (
	"io"

	"github.com/eswarantg/m3u8reader/common"
)

type M3u8Handler interface {
	PostRecord(tag common.TagId, kvpairs *AttrKVPairs) error
}

type Parser interface {
	ParseData(data []byte, handler M3u8Handler) (nBytes int, err error)
	Parse(rdr io.Reader, handler M3u8Handler) (nBytes int, err error)
}
