package m3u8reader

import (
	"bufio"
	"io"
)

type m3u8Handler interface {
	postRecord(tag string, kvpairs map[string]interface{}) error
}

const m3u8UnknownKey = "#"

func parseM3U8(src io.Reader, handler m3u8Handler) (nBytes int, err error) {

	brdr := bufio.NewReader(src)
	s := bufio.NewScanner(brdr)

	//Custom Split Function - Begin
	tokenCount := -1
	inTokenRead := false

	nBytes = 0
	lastTokenNewline := true

	custSplitFn := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		inQuotes := false
		tokenCount++
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
					if data[0] == '\n' {
						//String before
						nBytes += i
						return i, data[1:i], nil
					} else {
						//String before
						nBytes += i
						return i, data[0:i], nil
					}
				}
				nBytes += 1
				return 1, data[0:1], nil
			} else {
				if i > 0 && data[i-1] != '\n' {
					lastTokenNewline = false
				}
			}
			switch ch {
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
	kvpairs := make(map[string]interface{})
	lastToken := ""
	key := ""
	tag := ""
	postRecordFn := func() (err error) {
		if len(tag) > 0 {
			if len(key) > 0 {
				if _, ok := kvpairs[m3u8UnknownKey]; !ok {
					kvpairs[m3u8UnknownKey] = key
				} else {
					//Already present
					switch tag {
					case M3U8ExtInf:
						kvpairs["URI"] = key
					}
				}
				key = ""
			}
			//fmt.Printf("\npostRecordFn %v %v", tag, kvpairs)
			err = handler.postRecord(tag, kvpairs)
			tag = ""
			kvpairs = make(map[string]interface{})
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
					kvpairs[m3u8UnknownKey] = key
					key = ""
				}
				if curToken != "\n" {
					key = curToken
				}
			case "=":
				kvpairs[key] = curToken
				key = ""
			case "\n":
				if curToken != "\n" {
					if len(key) > 0 {
						kvpairs[key] = curToken
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
