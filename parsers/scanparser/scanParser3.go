package scanparser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type s3_ParsingState int

type ScanParser3 struct {
	extHandler parsers.M3u8Handler
	buffer     []byte

	//Parsing States
	state      s3_ParsingState
	savedState []s3_ParsingState
	tokenCount int
	nBytes     int

	kvpairs *parsers.AttrKVPairs
	tag     []byte
	tagId   common.TagId
	key     []byte
}

const (
	s3_UndefinedState s3_ParsingState = iota
	s3_ReadingQuote
	s3_ReadingEnumeratedString
	s3_ReadingAnyString
	s3_ReadingEntryName
	s3_WaitingEntryStart
	s3_WaitingEntryName
	s3_WaitingEntryData
)

func (s *ScanParser3) Init() {
	s.tokenCount = -1
	s.nBytes = 0
	s.state = s3_UndefinedState //bad state if we pop this
	s.savedState = make([]s3_ParsingState, 0, 10)
	s.pushState(s3_WaitingEntryStart) //push the starting state
	s.kvpairs = parsers.NewAttrKVPairs()
	s.tag = nil
	s.tagId = common.M3U8UNKNOWNTAG
	s.key = nil
}

func (s *ScanParser3) pushState(newState s3_ParsingState) {
	s.savedState = append(s.savedState, s.state)
	s.state = newState
}
func (s *ScanParser3) replaceState(newState s3_ParsingState) {
	//s.savedState = append(s.savedState, s.state)
	s.state = newState
}

func (s *ScanParser3) popState() {
	if len(s.savedState) <= 0 {
		panic("empty saved state stack")
	}
	s.state = s.savedState[len(s.savedState)-1]
	s.savedState = s.savedState[:len(s.savedState)-1]
}

func (s *ScanParser3) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
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

func (s *ScanParser3) Parse(rdr io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	s.extHandler = handler
	scan := bufio.NewScanner(rdr)
	if s.buffer == nil {
		s.buffer = make([]byte, 0, 4096)
	}
	scan.Buffer(s.buffer, len(s.buffer))
	return s.parse(scan, s)
}

func (s *ScanParser3) ParseData(data []byte, handler parsers.M3u8Handler) (nBytes int, err error) {
	s.extHandler = handler
	rdr := bytes.NewReader(data)
	scan := bufio.NewScanner(rdr)
	if s.buffer == nil {
		s.buffer = make([]byte, 0, 4096)
	}
	scan.Buffer(s.buffer, len(s.buffer))
	return s.parse(scan, s)
}

func (s *ScanParser3) postData(key common.AttrId, token []byte) error {
	if key == common.INTUnknownAttr && s.kvpairs.Exists(common.INTUnknownAttr) {
		switch s.tagId {
		case common.M3U8ExtInf:
			key = common.M3U8Uri
		default:
			return fmt.Errorf("duplicate INTUnknownAttr for %v required", s.tag)
		}
	}
	s.kvpairs.Store(key, string(token))
	return nil
}

func (s *ScanParser3) parse(scan *bufio.Scanner, handler parsers.M3u8Handler) (nBytes int, err error) {
	var lastToken []byte
	s.Init()
	scan.Split(s.splitFunctionMain)

	for scan.Scan() {
		s.tokenCount++
		curToken := scan.Bytes()
		//fmt.Printf("\nToken %v : %v : %v : %v", s.tokenCount, string(curToken), s.state, string(lastToken))
		if len(curToken) == 0 {
			continue
		}
		switch s.state {
		case s3_WaitingEntryStart:
			switch curToken[0] {
			case '#':
				s.popState()
				s.pushState(s3_WaitingEntryName)
				s.pushState(s3_ReadingEntryName)
			case '\n':
				//ignore empty line
			case ':', '=', ',':
				fallthrough
			default:
				err = fmt.Errorf("unexpected token %v received when waiting for EntryStart", string(curToken))
			}
		case s3_WaitingEntryName:
			switch curToken[0] {
			case ':':
				s.popState()
				s.pushState(s3_WaitingEntryData)
				s.pushState(s3_ReadingEnumeratedString)
				lastToken = nil
			case '\n':
				if s.tag != nil {
					s.PostRecord(s.tagId, s.kvpairs)
					s.tag = nil
					s.kvpairs = parsers.NewAttrKVPairs()
					s.popState()
					s.pushState(s3_WaitingEntryStart)
				}
				//ignore empty line
			case '#', '=', ',':
				err = fmt.Errorf(" %v unexpected token %v received when waiting for EntryName", string(s.tag), string(curToken))
			default:
				var ok bool
				s.tag = curToken
				s.tagId, ok = common.TagToTagId[string(curToken)]
				if !ok {
					err = fmt.Errorf("%v : invalid key token %v received when waiting for EntryName", string(s.tag), string(curToken))
				}
			}
		case s3_WaitingEntryData:
			switch curToken[0] {
			case '#':
				if s.tag != nil && lastToken[0] == '\n' {
					s.PostRecord(s.tagId, s.kvpairs)
					s.tag = nil
					s.kvpairs = parsers.NewAttrKVPairs()
					s.popState()
					s.pushState(s3_WaitingEntryName)
					s.pushState(s3_ReadingEntryName)
				} else {
					err = fmt.Errorf("%v : invalid token %v received when waiting for EntryData", string(s.tag), string(curToken))
				}
			case '=':
				if s.key != nil {
					err = fmt.Errorf("%v : key already realized %v, invalid token %v received when waiting for EntryData", string(s.tag), string(s.key), string(curToken))
				}
				s.key = lastToken
				lastToken = nil
				s.pushState(s3_ReadingAnyString)
			case '\n':
				if lastToken[0] == '\n' && s.tag != nil {
					s.PostRecord(s.tagId, s.kvpairs)
					s.tag = nil
					s.kvpairs = parsers.NewAttrKVPairs()
					s.popState()
					s.pushState(s3_WaitingEntryStart)
				}
				fallthrough
			case ',':
				if s.key != nil {
					var attr common.AttrId
					var ok bool
					attr, ok = common.AttrToAttrId[string(s.key)]
					if !ok {
						err = fmt.Errorf("invalid attribute token %v received when waiting for EntryData", string(s.key))
					}
					s.postData(attr, lastToken)
					s.key = nil
				} else {
					s.postData(common.INTUnknownAttr, lastToken)
				}
				lastToken = curToken
				s.pushState(s3_ReadingEnumeratedString)
			default:
				lastToken = curToken
			}
		}
		if err != nil {
			break
		}
	}
	if err == nil {
		if s.key != nil {
			var attr common.AttrId
			var ok bool
			attr, ok = common.AttrToAttrId[string(s.key)]
			if !ok {
				err = fmt.Errorf("invalid attribute token %v received when waiting for EntryData", string(s.key))
			}
			s.postData(attr, lastToken)
			s.key = nil
		} else {
			s.postData(common.INTUnknownAttr, lastToken)
		}
	}
	return s.nBytes, err
}

func (s *ScanParser3) splitFunctionMain(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch s.state {
	case s3_ReadingQuote:
		return s.readQuotedString(data, atEOF)
	case s3_ReadingEnumeratedString:
		return s.readEnumeratedString(data, atEOF)
	case s3_ReadingEntryName:
		return s.readEntryName(data, atEOF)
	case s3_ReadingAnyString:
		return s.readAnyString(data, atEOF)
	}
	for i, ch := range data {
		switch ch {
		case ' ':
			s.nBytes++
			continue //ignore space
		case '\r':
			s.nBytes++
			continue //ignore carriage return
		case '#', ':', '=', ',', '\n':
			s.nBytes++
			return i + 1, data[i : i+1], nil
		}
	}
	if atEOF {
		i := len(data)
		if i > 0 {
			s.popState()
			s.nBytes += i
			return i, data[0:i], nil
		}
	}
	return 0, nil, nil //need more characters
}

func (s *ScanParser3) readQuotedString(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i, ch := range data {
		switch ch {
		case '"':
			if i > 0 {
				s.popState()
				//read = i + 1 (skip quotes)
				s.nBytes += i + 1
				//data = 1:i  (skip last quote)
				return i, data[1:i], nil
			}
		case '\n', '\r':
			s.popState()
			return 0, nil, errors.New("multi-line quoted string not supported")
		}
	}
	if atEOF {
		s.popState()
		return 0, nil, errors.New("non-terminated quote string not supported")
	}
	return 0, nil, nil //need more characters
}
func (s *ScanParser3) readEntryName(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i, ch := range data {
		switch ch {
		case ' ':
			return 0, nil, errors.New("unexpected space reading entry name")
		case '\n', '\r', ':':
			s.popState()
			//don't include the delimiter
			//read = i (adjust for ZERO st value of i)...
			s.nBytes += i
			return i, data[0:i], nil
		}
	}
	if atEOF {
		i := len(data)
		if i > 0 {
			s.popState()
			s.nBytes += i
			return i, data[0:i], nil
		}
	}
	return 0, nil, nil //need more characters
}
func (s *ScanParser3) readEnumeratedString(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i, ch := range data {
		switch ch {
		case ' ':
			return 0, nil, errors.New("unexpected space reading enumerated string")
		case '\n', '\r', ',', '=', '#':
			s.popState()
			//don't include the delimiter
			//read = i (adjust for ZERO st value of i)...
			s.nBytes += i
			return i, data[0:i], nil
		}
	}
	if atEOF {
		i := len(data)
		if i > 0 {
			s.popState()
			s.nBytes += i
			return i, data[0:i], nil
		}
	}
	return 0, nil, nil //need more characters
}
func (s *ScanParser3) readAnyString(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if data[0] == '"' {
		s.replaceState(s3_ReadingQuote)
		return s.readQuotedString(data, atEOF)
	}
	s.replaceState(s3_ReadingEnumeratedString)
	return s.readEnumeratedString(data, atEOF)
}