package yaccparser

import (
	"fmt"
	"os"
	"testing"
)

type testParser struct {
	lval yySymType
}

func (p *testParser) Parse(l yyLexer) int {
	j := 0
	for {
		i := l.Lex(&p.lval)
		if i <= 0 {
			break
		}
		fmt.Printf("\n%v: %v %v %v", j, i, tokenIdToTagId(i), p.lval)
		j++
	}
	return 0
}
func (p *testParser) Lookahead() int {
	return 0
}

func Test_M3u8Lex(t *testing.T) {
	yyDebug = 5
	yyErrorVerbose = true
	tests := []string{
		//"test/ll_hls_byte_range.m3u8",
		//"test/ll_hls_delta_update.m3u8",
		//"test/ll_hls_pl.m3u8",
		//"test/index_new.m3u8",
		//"test/index_new_Variant_450k.m3u8",
		//"test/tv5.m3u8",
		//"test/tv5_TS-50002_1_video.m3u8",
		"test/master.m3u8",
		"test/sub.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		p := testParser{}
		p.Parse(NewLexer(f))
	}
}
