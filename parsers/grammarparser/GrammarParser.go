package grammarparser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type parserState int

const (
	searchingTag parserState = iota
	readingTag
	readingAttributes
	readingOpens
	readingNextline
)

type GrammarParser struct {
	state  parserState
	curTag common.TagId
	kv     parsers.AttrKVPairs
	line   int
	col    int
}

var boolToInt map[bool]int = map[bool]int{false: 0, true: 1}

func (p *GrammarParser) readFloat(data []byte, attrId common.AttrId) (value float64, remain []byte, err error) {
	//Assumption: value type is determined as valueQuotedString
	//checking of data[0] = '"' is already done, but it is sent in with data[0]
	pos := bytes.IndexAny(data, ",\n\r")
	if pos <= 0 {
		err = fmt.Errorf("line %v, Col %v : decimal float attribute %v value not found", p.line, p.col, common.AttrNames[attrId])
		return
	}
	value, err = strconv.ParseFloat(string(data[0:pos]), 64)
	if err != nil {
		err = fmt.Errorf("line %v, Col %v : decimal float attribute %v error parsing : %w", p.line, p.col, common.AttrNames[attrId], err)
		return
	}
	p.col += pos
	data = data[pos:]
	remain = data
	return
}

func (p *GrammarParser) readInt(data []byte, attrId common.AttrId) (value int64, remain []byte, err error) {
	//Assumption: value type is determined as valueQuotedString
	//checking of data[0] = '"' is already done, but it is sent in with data[0]
	pos := bytes.IndexAny(data, ",\n\r")
	if pos <= 0 {
		err = fmt.Errorf("line %v, Col %v : decimal int attribute %v value not found", p.line, p.col, common.AttrNames[attrId])
		return
	}
	value, err = strconv.ParseInt(string(data[0:pos]), 10, 64)
	if err != nil {
		err = fmt.Errorf("line %v, Col %v : decimal int attribute %v error parsing : %w", p.line, p.col, common.AttrNames[attrId], err)
		return
	}
	p.col += pos
	data = data[pos:]
	remain = data
	return
}

func (p *GrammarParser) readDateTime(data []byte, attrId common.AttrId) (value time.Time, remain []byte, err error) {
	//Assumption: value type is determined as valueQuotedString
	//checking of data[0] = '"' is already done, but it is sent in with data[0]
	pos := bytes.IndexAny(data, ",\n\r")
	if pos <= 0 {
		err = fmt.Errorf("line %v, Col %v : timeDate attribute %v value not found", p.line, p.col, common.AttrNames[attrId])
		return
	}
	value, err = time.Parse(time.RFC3339Nano, string(data[0:pos]))
	if err != nil {
		err = fmt.Errorf("line %v, Col %v : timeDate attribute %v error parsing : %w", p.line, p.col, common.AttrNames[attrId], err)
		return
	}
	p.col += pos
	data = data[pos:]
	remain = data
	return
}

func (p *GrammarParser) readQuotedString(data []byte, attrId common.AttrId) (value string, remain []byte, err error) {
	//Assumption: value type is determined as valueQuotedString
	//checking of data[0] = '"' is already done, but it is sent in with data[0]
	pos := bytes.IndexByte(data[1:], '"') //data[1:]- ignore first quote
	if pos < 0 {
		err = fmt.Errorf("line %v, Col %v : quoted string attribute %v quote not found", p.line, p.col, common.AttrNames[attrId])
		return
	}
	value = string(data[1 : pos+1]) //pos+1 - //adjust for the first quote
	pos += 2                        //adjust for the quotes
	p.col += pos
	data = data[pos:]
	remain = data
	return
}

func (p *GrammarParser) readEnumeratedString(data []byte, attrId common.AttrId) (value string, remain []byte, err error) {
	//Assumption: value type is determined as valueEnumeratedString
	pos := bytes.IndexAny(data, ",\n\r")
	if pos < 0 {
		err = fmt.Errorf("line %v, Col %v : enumerated string attribute %v value not found", p.line, p.col, common.AttrNames[attrId])
		return
	}
	value = string(data[0:pos])
	p.col += pos
	data = data[pos:]
	remain = data
	return
}

func (p *GrammarParser) readTag(data []byte) (remain []byte, err error) {
	//Assumption : must be after \n# point
	//Assumption : len(data)>0

	//find the position of next delimitter
	pos := bytes.IndexAny(data, ":\n")
	//Consume the entire content as tag if pos == -1 (not found)
	pos = boolToInt[pos == -1]*len(data) + boolToInt[pos != -1]*pos
	if pos == 0 {
		p.state = searchingTag
		return
	}
	tagId, ok := common.TagToTagId[string(data[0:pos])]
	if !ok {
		err = fmt.Errorf("line %v, Col %v : unknown tag \"%v\"", p.line, p.col, string(data[0:pos]))
		return
	}
	//setup the Tag
	p.curTag = tagId
	//create a new map for the attributes
	//Last record is owned by the post done
	p.kv = make(parsers.AttrKVPairs, 5)
	readBytes := pos + boolToInt[data[pos] == ':']
	p.col += readBytes - 1
	remain = data[readBytes:]
	meta := tagMeta[p.curTag]
	if meta.attrs != nil {
		p.state = readingAttributes
	} else if meta.openTypes != nil {
		p.state = readingOpens
	} else {
		p.state = searchingTag
	}
	return
}

func (p *GrammarParser) searchTag(data []byte) (remain []byte, err error) {
	//Assumption : must be before \n# point
	//Assumption : len(data)>0
	p.curTag = common.M3U8UNKNOWNTAG
	move1 := boolToInt[data[0] == '\r']
	move1 += boolToInt[data[move1] == '\n'] //move1 != 0 -  EOL read
	move2 := boolToInt[data[move1] == '\r']
	move2 += boolToInt[data[move1+move2] == '\n'] //move2 != 0 - another EOL read
	if move1 > 0 && move2 > 0 {                   //Empty line detected
		data = data[move1:]
		remain = data
		p.line++
		p.col = 0
		return
	}
	if move1 == 0 || data[move1] != '#' {
		err = fmt.Errorf("line %v, Col %v : invalid characters for new Tag", p.line, p.col)
		return
	}
	data = data[move1+1:]
	p.state = readingTag
	p.line++
	p.col = 1
	remain = data
	return
}

func (p *GrammarParser) getValue(data []byte, attrId common.AttrId, format ValueType) (value interface{}, remain []byte, err error) {
	//Assumption : must be after = for Attributes
	//Assumption : must be after : for Opens
	var valueStr string
	var floatVal float64
	var intVal int64
	var dateVal time.Time
	switch {
	case format&valueNextLineEnumeratedString > 0:
		if len(data) < 2 {
			err = fmt.Errorf("line %v, Col %v : next line enumerated string insufficient length", p.line, p.col)
			return
		}
		move := boolToInt[data[0] == '\r']
		move += boolToInt[data[move] == '\n']
		if move == 0 {
			err = fmt.Errorf("line %v, Col %v : next line enumerated string unable to find newline", p.line, p.col)
			return
		}
		data = data[move:]
		p.line++
		p.col = 0
		valueStr, data, err = p.readEnumeratedString(data, attrId)
		if err == nil {
			value = valueStr
		}
	case format&(valueQuotedString|valueEnumeratedString) > 0:
		if data[0] == '"' {
			valueStr, data, err = p.readQuotedString(data, attrId)
		} else {
			valueStr, data, err = p.readEnumeratedString(data, attrId)
		}
		if err == nil {
			value = valueStr
		}
	case format&valueQuotedString > 0:
		if data[0] != '"' {
			err = fmt.Errorf("line %v, Col %v : quoted string attribute %v quote not found", p.line, p.col, common.AttrNames[attrId])
			return
		}
		valueStr, data, err = p.readQuotedString(data, attrId)
		if err == nil {
			value = valueStr
		}
	case format&valueUTF8Text > 0:
		fallthrough //read as valueEnumeratedString for now
	case format&valueDecimalResolution > 0:
		fallthrough //treat same as valueEnumeratedString for now
	case format&valueEnumeratedString > 0:
		valueStr, data, err = p.readEnumeratedString(data, attrId)
		if err == nil {
			value = valueStr
		}
	case format&(valueSignedDecimalFloat|valueUnSignedDecimalFloat) > 0:
		floatVal, data, err = p.readFloat(data, attrId)
		if err == nil {
			value = floatVal
		}
	case format&valueDecimalInt > 0:
		intVal, data, err = p.readInt(data, attrId)
		if err == nil {
			value = intVal
		}
	case format&valueDateTime > 0:
		dateVal, data, err = p.readDateTime(data, attrId)
		if err == nil {
			value = dateVal
		}
	}
	remain = data
	return
}

func (p *GrammarParser) readingAttributes(data []byte) (remain []byte, err error) {
	//Assumption : must be after TAG:
	//Assumption : p.curTag is set to the current tag being read
	//Assumption : tagMeta[tagId].attrs is not NIL

	for {
		//Extract the Key
		//find the position of next delimitter
		pos := bytes.IndexAny(data, "=")
		if pos <= 0 {
			err = fmt.Errorf("line %v, Col %v : attribute tag not found", p.line, p.col)
			return
		}
		attrId, ok := common.AttrToAttrId[string(data[0:pos])]
		if !ok {
			err = fmt.Errorf("line %v, Col %v : unknown attribute tag %v", p.line, p.col, string(data[0:pos]))
			return
		}
		//Ignore checking if the attrId is supposed to be attribute for this Tag
		/*
			for _, item := range tagMeta[p.curTag].attrs {
				//Check the bitmap
				if attrId != item {
					err = fmt.Errorf("line %v, Col %v : attribute(%v) not be present in tag(%v)", p.line, p.col, common.AttrNames[attrId], tagNames[p.curTag])
				}
			}
		*/
		data = data[pos+1:] //ignoring the =
		p.col += pos + 1
		//Extract the value
		//find the position of next delimitter
		formats := attrMeta[attrId].types
		if formats == nil {
			err = fmt.Errorf("line %v, Col %v : unknown attribute %v value type not defined", p.line, p.col, common.AttrNames[attrId])
			return
		}
		var value interface{}
		value, data, err = p.getValue(data, attrId, formats[0])
		if err == nil {
			p.kv[attrId] = value
		}
		//position is at the seperator
		if data[0] != ',' {
			break
		}
		data = data[1:] //ignore ,
		p.col++
	}
	remain = data
	if tagMeta[p.curTag].openTypes != nil {
		p.state = readingOpens
	} else {
		p.state = searchingTag
	}
	return
}

func (p *GrammarParser) readingOpens(data []byte) (remain []byte, err error) {
	//Assumption : must be after TAG: or ATTRS are read
	//Assumption : p.curTag is set to the current tag being read
	//Assumption : tagMeta[tagId].openTypes is not NIL

	formats := tagMeta[p.curTag].openTypes
	var value interface{}
	//Extract the Value
	//find the position of next delimitter
	for _, format := range formats {
		value, data, err = p.getValue(data, format.attr, format.types)
		if err == nil {
			p.kv[format.attr] = value
		}
		// if one more item is present
		// Then move the data pointer
		if len(data) == 0 {
			break
		}
		if data[0] == ',' {
			data = data[1:]
			p.col++
		}
	}
	remain = data
	p.state = searchingTag
	return
}

func (p *GrammarParser) Parse(rdr io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	data, err := io.ReadAll(rdr)
	if err != nil {
		return len(data), err
	}
	return p.ParseData(data, handler)
}

func (p *GrammarParser) ParseData(data []byte, handler parsers.M3u8Handler) (nBytes int, err error) {
	origLen := len(data)
	p.state = readingTag
	p.kv = make(parsers.AttrKVPairs)
	if data[0] != '#' {
		err = fmt.Errorf("line %v, Col %v : expected # not found", p.line, p.col)
	}
	data = data[1:] //position after #
Loop:
	for err == nil {
		if len(data) <= 0 {
			break Loop
		}
		switch p.state {
		case readingTag:
			data, err = p.readTag(data)
		case searchingTag:
			if handler != nil && p.curTag != common.M3U8UNKNOWNTAG {
				err = handler.PostRecord(p.curTag, p.kv)
				if err != nil {
					break Loop
				}
			}
			if len(data) < 3 {
				break Loop
			}
			data, err = p.searchTag(data)
		case readingAttributes:
			data, err = p.readingAttributes(data)
		case readingOpens:
			data, err = p.readingOpens(data)
		}
	}
	return origLen - len(data), err
}
