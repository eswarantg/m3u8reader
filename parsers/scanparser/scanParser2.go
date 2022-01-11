package scanparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type ScanParser2 struct {
	extHander parsers.M3u8Handler
}

func (p *ScanParser2) SetBuffer([]byte) {
}

func (s *ScanParser2) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
	var err error
	if kvpairs != nil {
		err = decorateEntry(tag, *kvpairs)
		if err != nil {
			return err
		}
	}
	if s.extHander == nil {
		panic("\nInvalid extHandler for post")
	}
	err = s.extHander.PostRecord(tag, kvpairs)
	return err
}
func (s *ScanParser2) Parse(rdr io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	s.extHander = handler
	return parseM3U8_fast(rdr, s)
}

func (s *ScanParser2) ParseData(data []byte, handler parsers.M3u8Handler) (nBytes int, err error) {
	defer func() {
		s.extHander = nil
	}()
	s.extHander = handler
	rdr := bytes.NewReader(data)
	return parseM3U8_fast(rdr, s)
}

func stringSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the index of the input of a newline followed by a
	// pound sign.
	if i := strings.Index(string(data), "\n#"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}

func parseM3U8_fast(src io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	s := bufio.NewScanner(src)
	s.Split(stringSplitFunc)
	i := 0
	var kvpairs *parsers.AttrKVPairs
	kvpairs = parsers.NewAttrKVPairs() //initalize
	re := regexp.MustCompile(`[,\n]`)
	for s.Scan() {
		//fmt.Printf("\n%v", s.Text())
		entryStr := s.Text()
		index := strings.Index(entryStr, ":")
		tag := entryStr
		tagId := common.M3U8UNKNOWNTAG
		if index > 0 {
			tag = entryStr[:index]
			var ok bool
			tagId, ok = common.TagToTagId[tag]
			if !ok {
				panic(fmt.Sprintf("\nUnknown Tag : \"%v\"", tag))
			}
			split := re.Split(entryStr[index+1:], -1)
			for _, token := range split {
				if len(token) <= 0 {
					continue
				}
				parts := strings.Split(token, "=")
				if len(parts) == 2 {
					attrId, ok := common.AttrToAttrId[parts[0]]
					if ok {
						kvpairs.Store(attrId, parts[1])
					}
				} else {
					newkey := common.INTUnknownAttr
					if kvpairs != nil {
						if kvpairs.Exists(common.INTUnknownAttr) {
							//Already present
							switch tagId {
							case common.M3U8ExtInf:
								newkey = common.M3U8Uri
							default:
								//panic(fmt.Sprintf("Duplicate INTUnknownAttr for %v required.", tag))
							}
						}
					}
					kvpairs.Store(newkey, parts[0])
				}
				i++
			}
		}
		tagId, ok := common.TagToTagId[tag]
		if ok {
			handler.PostRecord(tagId, kvpairs)
			kvpairs = parsers.NewAttrKVPairs() //new value
		}
	}
	return 0, nil
}
