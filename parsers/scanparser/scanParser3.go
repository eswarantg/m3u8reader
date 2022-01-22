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

	//Parsing States
	state         s3_ParsingState
	savedState    [10]s3_ParsingState
	savedStateTop int
	tokenCount    int
	nBytes        int

	kvpairs *parsers.AttrKVPairs
	tag     []byte
	tagId   common.TagId
	key     []byte
	eof     bool
}

const (
	s3_UndefinedState s3_ParsingState = iota
	s3_ReadingQuote
	s3_ReadingEnumeratedString
	s3_ReadingAnyString
	s3_ReadingEntryName
	s3_ReadingIgnoredLine
	s3_WaitingEntryStart
	s3_WaitingEntryName
	s3_WaitingEntryData
)

func (s *ScanParser3) Init() {
	s.tokenCount = -1
	s.savedStateTop = -1
	s.nBytes = 0
	s.state = s3_UndefinedState       //bad state if we pop this
	s.pushState(s3_WaitingEntryStart) //push the starting state
	s.kvpairs = parsers.NewAttrKVPairs()
	s.tag = nil
	s.tagId = common.M3U8UNKNOWNTAG
	s.key = nil
	s.eof = false
}

func (s *ScanParser3) pushState(newState s3_ParsingState) {
	if s.savedStateTop >= len(s.savedState) {
		panic("saved state len is not enough")
	}
	s.savedStateTop++
	s.savedState[s.savedStateTop] = s.state
	s.state = newState
}
func (s *ScanParser3) replaceState(newState s3_ParsingState) {
	s.state = newState
}

func (s *ScanParser3) popState() {
	if s.savedStateTop < 0 {
		panic("empty saved state stack")
	}
	s.state = s.savedState[s.savedStateTop]
	s.savedStateTop--
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
		panic("Invalid extHandler for post")
	}
	err = s.extHandler.PostRecord(tag, kvpairs)
	return err
}

func (s *ScanParser3) Parse(rdr io.Reader, handler parsers.M3u8Handler, buffer []byte) (nBytes int, err error) {
	s.extHandler = handler
	scan := bufio.NewScanner(rdr)
	scan.Buffer(buffer, len(buffer))
	return s.parse(scan, s)
}

func (s *ScanParser3) ParseData(data []byte, handler parsers.M3u8Handler, buffer []byte) (nBytes int, err error) {
	s.extHandler = handler
	rdr := bytes.NewReader(data)
	scan := bufio.NewScanner(rdr)
	scan.Buffer(buffer, len(buffer))
	return s.parse(scan, s)
}

func (s *ScanParser3) postData(key common.AttrId, token []byte) error {
	if key == common.INTUnknownAttr {
		//fmt.Fprintf(os.Stdout, "\n%v:%v", string(s.tag), string(token))
		if s.kvpairs.Exists(common.INTUnknownAttr) {
			switch s.tagId {
			case common.M3U8ExtInf:
				key = common.M3U8Uri
			default:
				return fmt.Errorf("duplicate INTUnknownAttr for %v required", string(s.tag))
			}
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
			if !s.eof {
				continue
			}
			curToken = []byte{'\n'}
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
			//case ':', '=', ',':
			//fallthrough
			default:
				err = fmt.Errorf("unexpected token %v received when waiting for EntryStart", string(curToken))
			}
		case s3_WaitingEntryName:
			switch curToken[0] {
			case ':':
				if s.tag != nil {
					s.popState()
					s.pushState(s3_WaitingEntryData)
					s.pushState(s3_ReadingEnumeratedString)
					lastToken = nil
				}
			case '\n':
				if s.tag != nil {
					err = s.PostRecord(s.tagId, s.kvpairs)
					if err != nil {
						return s.nBytes, err
					}
				}
				s.tag = nil
				s.kvpairs = parsers.NewAttrKVPairs()
				s.popState()
				s.pushState(s3_WaitingEntryStart)
				lastToken = nil
				//ignore empty/commented line
			case '#', '=', ',':
				err = fmt.Errorf(" %v unexpected token %v received when waiting for EntryName", string(s.tag), string(curToken))
			default:
				var ok bool
				s.tag = curToken
				s.tagId, ok = common.TagToTagId[string(curToken)]
				if !ok {
					s.tag = nil
					s.pushState(s3_ReadingIgnoredLine)
					//ignore any line that is commented without valid key
				}
			}
		case s3_WaitingEntryData:
			switch curToken[0] {
			case '#':
				if s.tag != nil && lastToken[0] == '\n' {
					err = s.PostRecord(s.tagId, s.kvpairs)
					if err != nil {
						return s.nBytes, err
					}
					s.tag = nil
					s.kvpairs = parsers.NewAttrKVPairs()
					s.popState()
					s.pushState(s3_WaitingEntryName)
					s.pushState(s3_ReadingEntryName)
				} else {
					err = fmt.Errorf("%v : invalid token %v received when waiting for EntryData", string(s.tag), string(curToken))
					return s.nBytes, err
				}
			case '=':
				if s.key != nil {
					err = fmt.Errorf("%v : key already realized %v, invalid token %v received when waiting for EntryData", string(s.tag), string(s.key), string(curToken))
					return s.nBytes, err
				}
				s.key = lastToken
				lastToken = nil
				s.pushState(s3_ReadingAnyString)
			case '\n':
				//fmt.Fprintf(os.Stdout, "\nNEWLINE:%v", string(lastToken))
				if lastToken[0] == '\n' && s.tag != nil {
					err = s.PostRecord(s.tagId, s.kvpairs)
					if err != nil {
						return s.nBytes, err
					}
					s.tag = nil
					s.kvpairs = parsers.NewAttrKVPairs()
					s.popState()
					s.pushState(s3_WaitingEntryStart)
					continue
				}
				fallthrough
			case ',':
				//fmt.Fprintf(os.Stdout, "\nVALUE:%v", string(lastToken))
				if s.key != nil {
					var attr common.AttrId
					var ok bool
					attr, ok = common.AttrToAttrId[string(s.key)]
					if !ok {
						err = fmt.Errorf("invalid attribute token %v received when waiting for EntryData", string(s.key))
						return s.nBytes, err
					}
					err = s.postData(attr, lastToken)
					if err != nil {
						return s.nBytes, err
					}
					s.key = nil
					lastToken = nil
				} else {
					err = s.postData(common.INTUnknownAttr, lastToken)
					if err != nil {
						return s.nBytes, err
					}
					lastToken = nil
				}
				lastToken = curToken
				s.pushState(s3_ReadingEnumeratedString)
			default:
				lastToken = curToken
			}
		}
		if err != nil {
			return s.nBytes, err
		}
	}
	if scan.Err() != nil {
		err = scan.Err()
	}
	switch s.state {
	case s3_ReadingQuote:
		fallthrough
	case s3_ReadingEnumeratedString:
		s.popState()
	}
	switch s.state {
	case s3_WaitingEntryStart:
		//fmt.Printf("%v %v %v", "s3_WaitingEntryStart", s.tag, lastToken)
	case s3_WaitingEntryName:
		//fmt.Printf("%v %v %v", "s3_WaitingEntryName", s.tag, lastToken)
		if len(lastToken) > 0 {
			var ok bool
			s.tag = lastToken
			s.tagId, ok = common.TagToTagId[string(lastToken)]
			if !ok {
				s.tag = nil
				//ignore the line ... might be some comment line
			}
		}
		if len(s.tag) > 0 {
			err = s.PostRecord(s.tagId, s.kvpairs)
			if err != nil {
				return s.nBytes, err
			}
		}
	case s3_WaitingEntryData:
		//fmt.Printf("%v %v %v", "s3_WaitingEntryData", string(s.tag), string(lastToken))
		if len(lastToken) > 0 && lastToken[0] != '\n' {
			if s.key != nil {
				var attr common.AttrId
				var ok bool
				attr, ok = common.AttrToAttrId[string(s.key)]
				if !ok {
					err = fmt.Errorf("invalid attribute token %v received when waiting for EntryData", string(s.key))
					return s.nBytes, err
				}
				err = s.postData(attr, lastToken)
				if err != nil {
					return s.nBytes, err
				}
				lastToken = nil
				s.key = nil
			} else {
				err = s.postData(common.INTUnknownAttr, lastToken)
				if err != nil {
					return s.nBytes, err
				}
				lastToken = nil
			}
		}
		if len(s.tag) > 0 {
			err = s.PostRecord(s.tagId, s.kvpairs)
			if err != nil {
				return s.nBytes, err
			}
		}
	default:
		panic(fmt.Sprintf("Unexpected state : %v", s.state))
	}
	if err == nil {
		if s.key != nil {
			var attr common.AttrId
			var ok bool
			attr, ok = common.AttrToAttrId[string(s.key)]
			if !ok {
				err = fmt.Errorf("invalid attribute token %v received when waiting for EntryData", string(s.key))
				return s.nBytes, err
			}
			err = s.postData(attr, lastToken)
			if err != nil {
				return s.nBytes, err
			}
			lastToken = nil
			s.key = nil
		} else {
			if len(lastToken) > 0 && lastToken[0] != '\n' {
				err = s.postData(common.INTUnknownAttr, lastToken)
				if err != nil {
					return s.nBytes, err
				}
				lastToken = nil
			}
		}
	}
	return s.nBytes, err
}

func (s *ScanParser3) splitFunctionMain(data []byte, atEOF bool) (int, []byte, error) {
	switch s.state {
	case s3_ReadingQuote:
		return s.readQuotedString(data, atEOF)
	case s3_ReadingEnumeratedString:
		return s.readEnumeratedString(data, atEOF)
	case s3_ReadingEntryName:
		return s.readEntryName(data, atEOF)
	case s3_ReadingAnyString:
		return s.readAnyString(data, atEOF)
	case s3_ReadingIgnoredLine:
		return s.readIgnoredLine(data, atEOF)
	}
	for i, ch := range data {
		if ch == '#' || ch == ':' || ch == '=' || ch == ',' || ch == '\n' {
			s.nBytes += i
			return i + 1, data[i : i+1], nil
		}
	}
	if atEOF {
		i := len(data)
		if i > 0 {
			s.nBytes += i
			s.popState()
			return i, data[0:i], nil
		}
		s.eof = true
		return 0, nil, io.EOF
	}
	return 0, nil, nil //need more characters
}

func (s *ScanParser3) readQuotedString(data []byte, atEOF bool) (int, []byte, error) {
	for i, ch := range data {
		if ch == '"' {
			if i > 0 {
				s.nBytes += i
				s.popState()
				//data = 1:i  (skip last quote)
				return i, data[1:i], nil
			}
		} else if ch == '\r' || ch == '\n' {
			s.nBytes += i
			s.popState()
			return 0, nil, errors.New("multi-line quoted string not supported")
		}
	}
	if atEOF {
		s.eof = true
		s.nBytes += len(data)
		s.popState()
		return 0, nil, errors.New("non-terminated quote string not supported")
	}
	return 0, nil, nil //need more characters
}
func (s *ScanParser3) readEntryName(data []byte, atEOF bool) (int, []byte, error) {
	for i, ch := range data {
		if ch == '\r' || ch == '\n' || ch == ':' {
			s.nBytes += i
			s.popState()
			return i, data[0:i], nil
		}
	}
	if atEOF {
		i := len(data)
		if i > 0 {
			s.nBytes += i
			s.popState()
			return i, data[0:i], nil
		}
		s.eof = true
		return 0, nil, io.EOF
	}
	return 0, nil, nil //need more characters
}

func (s *ScanParser3) readIgnoredLine(data []byte, atEOF bool) (int, []byte, error) {
	for i, ch := range data {
		if ch == '\r' || ch == '\n' {
			s.nBytes += i
			s.popState()
			return i, data[0:i], nil
		}
	}
	if atEOF {
		i := len(data)
		if i > 0 {
			s.nBytes += i
			s.popState()
			return i, data[0:i], nil
		}
		s.eof = true
		return 0, nil, io.EOF
	}
	return 0, nil, nil //need more characters
}
func (s *ScanParser3) readEnumeratedString(data []byte, atEOF bool) (int, []byte, error) {
	for i, ch := range data {
		if ch == '\r' || ch == '\n' || ch == ',' || ch == '=' || ch == '#' {
			s.nBytes += i
			s.popState()
			//don't include the delimiter
			//read = i (adjust for ZERO st value of i)...
			return i, data[0:i], nil
		} else if ch == ' ' || ch == '"' {
			s.nBytes += i
			s.popState()
			return 0, nil, fmt.Errorf("unexpected char (%v) reading enumerated string", ch)
		}
	}
	if atEOF {
		s.eof = true
		i := len(data)
		if i > 0 {
			s.nBytes += i
			s.popState()
			return i, data[0:i], nil
		}
		s.eof = true
		return 0, nil, io.EOF
	}
	return 0, nil, nil //need more characters
}
func (s *ScanParser3) readAnyString(data []byte, atEOF bool) (int, []byte, error) {
	if data[0] == '"' {
		s.replaceState(s3_ReadingQuote)
		return s.readQuotedString(data, atEOF)
	}
	s.replaceState(s3_ReadingEnumeratedString)
	return s.readEnumeratedString(data, atEOF)
}
