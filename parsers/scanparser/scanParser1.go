package scanparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type ScanParser1 struct {
	extHandler parsers.M3u8Handler
	buffer     []byte
}

func (s *ScanParser1) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
	var err error
	if kvpairs != nil {
		err = decorateEntry(tag, *kvpairs)
		if err != nil {
			return err
		}
	}
	if s.extHandler == nil {
		panic("\nInvalid extHandler for post")
	}
	err = s.extHandler.PostRecord(tag, kvpairs)
	return err
}

func (s *ScanParser1) Parse(rdr io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	s.extHandler = handler
	scan := bufio.NewScanner(rdr)
	if s.buffer == nil {
		s.buffer = make([]byte, 0, 4096)
	}
	scan.Buffer(s.buffer, len(s.buffer))
	return parseM3U8(scan, s)
}

func (s *ScanParser1) ParseData(data []byte, handler parsers.M3u8Handler) (nBytes int, err error) {
	s.extHandler = handler
	rdr := bytes.NewReader(data)
	scan := bufio.NewScanner(rdr)
	if s.buffer == nil {
		s.buffer = make([]byte, 0, 4096)
	}
	scan.Buffer(s.buffer, len(s.buffer))
	return parseM3U8(scan, s)
}

func parseM3U8(s *bufio.Scanner, handler parsers.M3u8Handler) (nBytes int, err error) {
	//Custom Split Function - Begin
	tokenCount := -1
	inTokenRead := false

	nBytes = 0
	lastTokenNewline := true

	custSplitFn := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		inQuotes := false
		tokenCount++
		carriageReturnRead := false
		for i, ch := range data {
			if inTokenRead {
				if ch != ':' && ch != '\n' {
					continue
				}
				if i == 0 {
					inTokenRead = false
					if ch == '\n' {
						lastTokenNewline = true
						//Skip new line
						nBytes += 1
						return 1, data[0:1], nil
					} else {
						//Token ':'
						nBytes += 1
						return 1, data[0:1], nil
					}
				} else {
					//String before
					nBytes += i
					return i, data[0:i], nil
				}
			}
			if inQuotes {
				//if ch == '\n' {
				//TBD
				//}
				if ch != '"' {
					continue
				}
				inQuotes = false
				//String before in quotes
				nBytes += i + 1
				return i + 1, data[1:i], nil
			}
			if ch == '\n' {
				lastTokenNewline = true
				if i > 0 {
					if data[0] == '\n' || data[0] == '\r' {
						//String before
						nBytes += i
						if carriageReturnRead {
							return i, data[2:i], nil
						} else {
							return i, data[1:i], nil
						}
					} else {
						//String before
						nBytes += i
						return i, data[0:i], nil
					}
				}
				if carriageReturnRead {
					nBytes += 2
					return 2, data[1:2], nil
				} else {
					nBytes += 1
					return 1, data[0:1], nil
				}
			} else {
				if i > 0 && data[i-1] != '\n' {
					lastTokenNewline = false
				}
				carriageReturnRead = false
			}
			switch ch {
			case '\r':
				carriageReturnRead = true
				continue
			case '"':
				inQuotes = true
				continue
			case '#':
				if (i == 1 || i == 0) && lastTokenNewline {
					inTokenRead = true
					//Token
					nBytes += i + 1
					return i + 1, data[i : i+1], nil
				}
			case ',', '=':
				if i == 0 {
					//Token
					nBytes += 1
					return 1, data[0 : i+1], nil
				}
				//String before
				nBytes += i
				return i, data[0:i], nil
			}
		}
		if atEOF && len(data) > 0 {
			nBytes += len(data)
			return len(data), data, nil
		}
		return 0, nil, nil
	}

	s.Split(custSplitFn)
	//Custom Split Function - End

	//Post Record Entry - Start
	var kvpairs *parsers.AttrKVPairs
	kvpairs = parsers.NewAttrKVPairs() //initalize
	var lastToken []byte
	var key []byte
	var tag []byte
	tagId := common.M3U8UNKNOWNTAG
	postRecordFn := func() (err error) {
		if tag != nil {
			if key != nil {
				newkey := common.INTUnknownAttr
				if kvpairs.Exists(common.INTUnknownAttr) {
					//Already present
					switch tagId {
					case common.M3U8ExtInf:
						newkey = common.M3U8Uri
					default:
						//panic(fmt.Sprintf("Duplicate INTUnknownAttr for %v required.", tag))
					}
				}
				kvpairs.Store(newkey, string(key))
				key = nil
			}
			//fmt.Printf("\npostRecordFn %v %v", tag, kvpairs)
			tagId, ok := common.TagToTagId[string(tag)]
			if ok {
				err = handler.PostRecord(tagId, kvpairs)
				kvpairs = parsers.NewAttrKVPairs() //use new one next time
			}
			tag = nil
		}
		return
	}
	//Post Record Entry - End

	for s.Scan() {
		//fmt.Printf("\nToken %v : %v", tokenCount, s.Text())
		curToken := s.Bytes()
		if len(curToken) == 0 {
			continue
		}
		if bytes.Equal(curToken, []byte{'#'}) {
			err = postRecordFn()
			if err != nil {
				break
			}
		} else {
			switch lastToken[0] {
			case '#':
				if curToken[0] == '\n' {
					continue //skip new line with only #
				}
				var ok bool
				tag = curToken
				tagId, ok = common.TagToTagId[string(tag)]
				if !ok {
					panic(fmt.Sprintf("\nUnknown Tag : \"%v\"", string(tag)))
				}
			case ',', ':':
				if key != nil {
					kvpairs.Store(common.INTUnknownAttr, string(key))
					key = nil
				}
				if curToken[0] != '\n' {
					key = curToken
				}
			case '=':
				attrId, ok := common.AttrToAttrId[string(key)]
				if ok {
					kvpairs.Store(attrId, string(curToken))
				}
				key = nil
			case '\n':
				if curToken[0] != '\n' {
					if key != nil {
						attrId, ok := common.AttrToAttrId[string(key)]
						if ok {
							kvpairs.Store(attrId, string(curToken))
						}
						key = nil
					} else {
						key = curToken
					}
				}
			}
		}
		lastToken = s.Bytes()
	}
	if err == nil {
		err = postRecordFn()
	}
	return
}
