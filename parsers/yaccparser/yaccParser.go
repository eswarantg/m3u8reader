package yaccparser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

type YaccParser struct {
	extHander parsers.M3u8Handler
}

func (y *YaccParser) PostRecord(tag common.TagId, kvpairs parsers.AttrKVPairs) error {
	return y.extHander.PostRecord(tag, kvpairs)
}

func (y *YaccParser) yyparse(r io.Reader, handler parsers.M3u8Handler) (nbytes int, err error) {
	nbytes = -1 //don't know how to get number of bytes read... for now put -1
	lex := NewLexerWithInit(r, func(l *Lexer) {
		l.parseResult = handler
	})
	defer func() {
		if err1 := recover(); err1 != nil {
			msg, ok := err1.(string)
			if ok {
				err = fmt.Errorf("%v : at line %v", msg, lex.Line())
				return
			}
			err = fmt.Errorf("%v : panic handled : at line %v", err1, lex.Line())
			return
		}
	}()
	result := yyParse(lex)
	if result != 0 && err == nil {
		err = fmt.Errorf("yyParse returned non-zero")
	}
	return
}

func (y *YaccParser) Parse(rdr io.Reader, handler parsers.M3u8Handler) (nBytes int, err error) {
	defer func() {
		y.extHander = nil
	}()
	y.extHander = handler
	return y.yyparse(rdr, handler)
}

func (y *YaccParser) ParseData(data []byte, handler parsers.M3u8Handler) (nBytes int, err error) {
	defer func() {
		y.extHander = nil
	}()
	y.extHander = handler
	rdr := bytes.NewReader(data)
	return y.yyparse(rdr, handler)
}
