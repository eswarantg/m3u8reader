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
	extHander parsers.M3u8Handler
}

func (s *ScanParser1) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
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

func (s *ScanParser1) Parse(rdr io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	s.extHander = handler
	return parseM3U8(rdr, s)
}

func (s *ScanParser1) ParseData(data []byte, handler parsers.M3u8Handler) (nBytes int, err error) {
	defer func() {
		s.extHander = nil
	}()
	s.extHander = handler
	rdr := bytes.NewReader(data)
	return parseM3U8(rdr, s)
}

func parseM3U8(src io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {

	s := bufio.NewScanner(src)

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
	lastToken := ""
	key := ""
	tag := ""
	tagId := common.M3U8UNKNOWNTAG
	postRecordFn := func() (err error) {
		if len(tag) > 0 {
			if len(key) > 0 {
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
				kvpairs.Store(newkey, key)
				key = ""
			}
			//fmt.Printf("\npostRecordFn %v %v", tag, kvpairs)
			tagId, ok := common.TagToTagId[tag]
			if ok {
				err = handler.PostRecord(tagId, kvpairs)
				kvpairs = parsers.NewAttrKVPairs() //use new one next time
			}
			tag = ""
		}
		return
	}
	//Post Record Entry - End

	for s.Scan() {
		//fmt.Printf("\nToken %v : %v", tokenCount, s.Text())
		curToken := s.Text()
		if len(curToken) == 0 {
			continue
		}
		if curToken == "#" {
			err = postRecordFn()
			if err != nil {
				break
			}
		} else {
			switch lastToken {
			case "#":
				if curToken == "\n" {
					continue //skip new line with only #
				}
				var ok bool
				tag = curToken
				tagId, ok = common.TagToTagId[tag]
				if !ok {
					panic(fmt.Sprintf("\nUnknown Tag : \"%v\"", tag))
				}
			case ",", ":":
				if len(key) > 0 {
					kvpairs.Store(common.INTUnknownAttr, key)
					key = ""
				}
				if curToken != "\n" {
					key = curToken
				}
			case "=":
				attrId, ok := common.AttrToAttrId[key]
				if ok {
					kvpairs.Store(attrId, curToken)
				}
				key = ""
			case "\n":
				if curToken != "\n" {
					if len(key) > 0 {
						attrId, ok := common.AttrToAttrId[key]
						if ok {
							kvpairs.Store(attrId, curToken)
						}
						key = ""
					} else {
						key = curToken
					}
				}
			}
		}
		lastToken = s.Text()
	}
	if err == nil {
		err = postRecordFn()
	}
	return
}
