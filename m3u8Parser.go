package m3u8reader

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type m3u8Handler interface {
	postRecord(tag TagId, kvpairs map[AttrId]interface{}) error
}

const m3u8UnknownKey = "#"

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

func parseM3U8_fast(src io.Reader, handler m3u8Handler) (nBytes int, err error) {
	s := bufio.NewScanner(src)
	s.Split(stringSplitFunc)
	i := 0
	kvpairs := make(map[AttrId]interface{}, 4)
	re := regexp.MustCompile(`[,\n]`)
	for s.Scan() {
		//fmt.Printf("\n%v", s.Text())
		entryStr := s.Text()
		index := strings.Index(entryStr, ":")
		tag := entryStr
		if index > 0 {
			tag = entryStr[:index]
			split := re.Split(entryStr[index+1:], -1)
			for _, token := range split {
				if len(token) <= 0 {
					continue
				}
				parts := strings.Split(token, "=")
				if len(parts) == 2 {
					attrId, ok := attrToAttrId[parts[0]]
					if ok {
						kvpairs[attrId] = parts[1]
					}
				} else {
					kvpairs[INTUnknownAttr] = parts[0]
				}
				i++
			}
		}
		tagId, ok := tagToTagId[tag]
		if ok {
			handler.postRecord(tagId, kvpairs)
		}
		for k := range kvpairs {
			delete(kvpairs, k)
		}
	}
	return 0, nil
}

func parseM3U8(src io.Reader, handler m3u8Handler) (nBytes int, err error) {

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
	kvpairs := make(map[AttrId]interface{}, 5)
	lastToken := ""
	key := ""
	tag := ""
	postRecordFn := func() (err error) {
		if len(tag) > 0 {
			if len(key) > 0 {
				if _, ok := kvpairs[INTUnknownAttr]; !ok {
					kvpairs[INTUnknownAttr] = key
				} else {
					//Already present
					switch tagToTagId[tag] {
					case M3U8ExtInf:
						kvpairs[M3U8Uri] = key
					}
				}
				key = ""
			}
			//fmt.Printf("\npostRecordFn %v %v", tag, kvpairs)
			tagId, ok := tagToTagId[tag]
			if ok {
				err = handler.postRecord(tagId, kvpairs)
			}
			tag = ""
			kvpairs = make(map[AttrId]interface{}, 5)
		}
		return
	}
	//Post Record Entry - End

	for s.Scan() {
		//fmt.Printf("\nToken %v : %v", tokenCount, s.Text())
		curToken := s.Text()
		if curToken == "#" {
			err = postRecordFn()
			if err != nil {
				break
			}
		} else {
			switch lastToken {
			case "#":
				tag = curToken
			case ",", ":":
				if len(key) > 0 {
					kvpairs[INTUnknownAttr] = key
					key = ""
				}
				if curToken != "\n" {
					key = curToken
				}
			case "=":
				attrId, ok := attrToAttrId[key]
				if ok {
					kvpairs[attrId] = curToken
				}
				key = ""
			case "\n":
				if curToken != "\n" {
					if len(key) > 0 {
						attrId, ok := attrToAttrId[key]
						if ok {
							kvpairs[attrId] = curToken
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
