package grammarparser

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/eswarantg/m3u8reader/common"
)

func Test_readFloat(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{"1.000,", "-2.000\nabcd", "1\r\n", "2.0abcd,", "abcd", "abcd,", ",abcd", ""}
	results := [...]float64{1.000, -2.000, 1.000, 0, 0, 0, 0, 0}
	remains := [...]int{1, 5, 2, 0, 0, 0, 0, 0}
	cols := [...]int{5, 6, 1, 0, 0, 0, 0, 0}
	errors := [...]error{nil, nil, nil, errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR")}
	for i, sample := range samples {
		p.col = 0
		val, remain, err := p.readFloat([]byte(sample), common.INTUnknownAttr)
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if remains[i] != len(remain) {
			t.Errorf("remain expected %v : got %v", remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("no error expected %v : got %v", errors[i], err.Error())
		}
		if results[i] != val {
			t.Errorf("parsed value not match expected %v : got %v", results[i], val)
		}
		if cols[i] != p.col {
			t.Errorf("cols expected %v : got %v", cols[i], p.col)
		}
	}
}

func Test_readInt(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{"1234,", "-2234\nabcd", "1\r\n", "2.0,abcd,", "abcd", "abcd,", ",abcd", ""}
	results := [...]int64{1234, -2234, 1, 0, 0, 0, 0, 0}
	remains := [...]int{1, 5, 2, 0, 0, 0, 0, 0}
	cols := [...]int{4, 5, 1, 0, 0, 0, 0, 0}
	errors := [...]error{nil, nil, nil, errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR")}
	for i, sample := range samples {
		p.col = 0
		val, remain, err := p.readInt([]byte(sample), common.INTUnknownAttr)
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : no error expected %v : got %v", i, errors[i], err.Error())
		}
		if results[i] != val {
			t.Errorf("%v : parsed value not match expected %v : got %v", i, results[i], val)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
	}
}

func Test_readEnumeratedString(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{"1234,", "-2234\nabcd", "1\r\n", "2.0,abcd,", "abcd", "abcd,", ",abcd", ""}
	results := [...]string{"1234", "-2234", "1", "2.0", "", "abcd", "", ""}
	remains := [...]int{1, 5, 2, 6, 0, 1, 5, 0}
	cols := [...]int{4, 5, 1, 3, 0, 4, 0, 0}
	errors := [...]error{nil, nil, nil, nil, errors.New("ERR"), nil, nil, errors.New("ERR")}
	for i, sample := range samples {
		p.col = 0
		val, remain, err := p.readEnumeratedString([]byte(sample), common.INTUnknownAttr)
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		if results[i] != val {
			t.Errorf("%v : parsed value not match expected \"%v\" : got \"%v\"", i, results[i], val)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
	}
}

func Test_readQuotedString(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{"\"1234\",", "\"-2234\"\nabcd", "\"1\"\r\n", "\"2.0\",abcd,", "\"abcd", "\"abcd\",", "\",abcd", "\"\"", "\""}
	results := [...]string{"1234", "-2234", "1", "2.0", "", "abcd", "", "", ""}
	remains := [...]int{1, 5, 2, 6, 0, 1, 0, 0, 0}
	cols := [...]int{6, 7, 3, 5, 0, 6, 0, 2, 0}
	errors := [...]error{nil, nil, nil, nil, errors.New("ERR"), nil, errors.New("ERR"), nil, errors.New("ERR")}
	for i, sample := range samples {
		p.col = 0
		val, remain, err := p.readQuotedString([]byte(sample), common.INTUnknownAttr)
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		if results[i] != val {
			t.Errorf("%v : parsed value not match expected \"%v\" : got \"%v\"", i, results[i], val)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
	}
}

func Test_readDateTime(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{"2021-11-29T22:24:07.311Z\n", "2021-11-29T22:24:07.311Z,abcd", "2021-11-29T22:24:07.311Z\r\n#", "2021-11-29T22:24:07.311Z"}
	results := [...]string{"2021-11-29T22:24:07.311Z", "2021-11-29T22:24:07.311Z", "2021-11-29T22:24:07.311Z", "0001-01-01T00:00:00Z"}
	remains := [...]int{1, 5, 3, 0}
	cols := [...]int{24, 24, 24, 0}
	errors := [...]error{nil, nil, nil, errors.New("ERR")}
	for i, sample := range samples {
		p.col = 0
		val, remain, err := p.readDateTime([]byte(sample), common.INTUnknownAttr)
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		fmtVal := val.Format(time.RFC3339Nano)
		if results[i] != fmtVal {
			t.Errorf("%v : parsed value not match expected \"%v\" : got \"%v\"", i, results[i], fmtVal)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
	}
}

func Test_readTag(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{
		"EXTM3U\n",
		"EXT-X-VERSION:9\n",
		"EXT-X-TARGETDURATION:4\n",
		"EXT-X-SERVER-CONTROL:CAN-BLOCK-RELOAD=YES,CAN-SKIP-UNTIL=24,PART-HOLD-BACK=3.012\n",
		"EXTINF:4.00000,\nfileSequence436248.m4s\n",
		"EXT-X-STREAM-INF:BANDWIDTH=550172,RESOLUTION=256x106\nlevel_0.m3u8\n",
		"EXT-X-PART:DURATION=1.000,URI=\"tv5_TS-50002_1_video_91001847.0.mp4\",INDEPENDENT=YES\n",
		"EXT-X-PART-INF:PART-TARGET=1.004000\n",
		":abcd\n",
		"abcd\n",
		"\n",
		":abcd",
		":",
		"",
	}
	remains := [...]int{1, 2, 2, 60, 32, 49, 73, 21, 0, 0, 0, 0, 0, 0}
	tags := [...]common.TagId{common.M3U8FormatIdentifier, common.M3U8ExtXVersion, common.M3U8TargetDuration, common.M3U8ExtXServerControl, common.M3U8ExtInf, common.M3U8ExtXStreamInf, common.M3U8ExtXPart, common.M3U8ExtXPartInf,
		common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG}
	states := [...]parserState{searchingTag, readingOpens, readingOpens, readingAttributes, readingOpens, readingAttributes, readingAttributes, readingAttributes,
		readingTag, readingTag, readingTag, readingTag, readingTag, readingTag}
	cols := [...]int{5, 13, 20, 20, 6, 16, 10, 14,
		0, 0, 0, 0, 0, 0}
	errors := [...]error{nil, nil, nil, nil, nil, nil, nil, nil,
		errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR")}
	for i, sample := range samples {
		p.state = readingTag
		p.curTag = common.M3U8UNKNOWNTAG
		p.line = 0
		p.col = 0
		remain, err := p.readTag([]byte(sample))
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if tags[i] != p.curTag {
			t.Errorf("%v : Tag expected %v : got %v", i, tags[i], p.curTag)
		}
		if states[i] != p.state {
			t.Errorf("%v : state expected %v : got %v", i, states[i], p.state)
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
	}
}

func Test_searchTag(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{
		"\r\n#EXT-X-VERSION:9\n",      //0
		"\n#EXT-X-TARGETDURATION:4\n", //1
		"\n#EXT-X-SERVER-CONTROL:CAN-BLOCK-RELOAD=YES,CAN-SKIP-UNTIL=24,PART-HOLD-BACK=3.012\n",      //2
		"\n#EXTINF:4.00000,\nfileSequence436248.m4s\n",                                               //3
		"\n#EXT-X-STREAM-INF:BANDWIDTH=550172,RESOLUTION=256x106\nlevel_0.m3u8\n",                    //4
		"\n#EXT-X-PART:DURATION=1.000,URI=\"tv5_TS-50002_1_video_91001847.0.mp4\",INDEPENDENT=YES\n", //5
		"\n#EXT-X-PART-INF:PART-TARGET=1.004000\n",                                                   //6
		"\r\n#", //7
		"\n#E",  //8
		"abcd",  //9
		"\r\n",  //10
		"\n#",   //11
		"\n",    //12
	}
	remains := [...]int{16, 23, 81, 39, 66, 84, 36,
		0, 1, 0, 0, 0, 0}
	tags := [...]common.TagId{common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG,
		common.M3U8UNKNOWNTAG, common.M3U8UNKNOWNTAG, common.M3U8FormatIdentifier, common.M3U8FormatIdentifier, common.M3U8FormatIdentifier, common.M3U8FormatIdentifier}
	states := [...]parserState{readingTag, readingTag, readingTag, readingTag, readingTag, readingTag, readingTag,
		readingTag, readingTag, searchingTag, searchingTag, searchingTag, searchingTag}
	cols := [...]int{1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 0}
	lines := [...]int{1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 0}
	errors := [...]error{nil, nil, nil, nil, nil, nil, nil,
		nil, nil, errors.New("ERR"), errors.New("ERR"), errors.New("ERR"), errors.New("ERR")}
	for i, sample := range samples {
		p.state = searchingTag
		p.curTag = common.M3U8FormatIdentifier //Sample
		p.line = 0
		p.col = 0
		remain, err := p.searchTag([]byte(sample))
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if tags[i] != p.curTag {
			t.Errorf("%v : Tag expected %v : got %v", i, tags[i], p.curTag)
		}
		if states[i] != p.state {
			t.Errorf("%v : state expected %v : got %v", i, states[i], p.state)
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
		if lines[i] != p.line {
			t.Errorf("%v : lines expected %v : got %v", i, lines[i], p.line)
		}
	}
}

func Test_readAttributes(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{
		"CAN-BLOCK-RELOAD=YES,CAN-SKIP-UNTIL=24,PART-HOLD-BACK=3.012\n",
		"BANDWIDTH=550172,RESOLUTION=256x106\nlevel_0.m3u8\n",
		"DURATION=1.000,URI=\"tv5_TS-50002_1_video_91001847.0.mp4\",INDEPENDENT=YES\n",
		"PART-TARGET=1.004000\n",
	}
	kvs := [...]map[common.AttrId]interface{}{
		{common.M3U8CanBlockReload: "YES", common.M3U8CanSkipUntil: int64(24), common.M3U8PartHoldBack: float64(3.012)},
		{common.M3U8Bandwidth: int64(550172), common.M3U8Resolution: "256x106"},
		{common.M3U8Duration: float64(1.000), common.M3U8Uri: "tv5_TS-50002_1_video_91001847.0.mp4", common.M3U8Independent: "YES"},
		{common.M3U8PartTarget: float64(1.004000)},
	}
	tags := [...]common.TagId{common.M3U8ExtXServerControl, common.M3U8ExtXStreamInf, common.M3U8ExtXPart, common.M3U8ExtXPartInf}
	remains := [...]int{1, 14, 1, 1}
	states := [...]parserState{searchingTag, readingOpens, searchingTag, searchingTag}
	cols := [...]int{59, 35, 72, 20}
	lines := [...]int{0, 0, 0, 0}
	errors := [...]error{nil, nil, nil, nil,
		errors.New("ERR")}
	for i, sample := range samples {
		p.state = searchingTag
		p.curTag = tags[i]
		p.line = 0
		p.col = 0
		p.kv = nil
		remain, err := p.readingAttributes([]byte(sample))
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if states[i] != p.state {
			t.Errorf("%v : state expected %v : got %v", i, states[i], p.state)
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
		if lines[i] != p.line {
			t.Errorf("%v : lines expected %v : got %v", i, lines[i], p.line)
		}
		for k, v := range kvs[i] {
			var val interface{}
			var ok bool
			if val = p.kv.Get(k); val != nil {
				t.Errorf("%v : kv expected %v : not found", i, k)
			}
			switch ty := reflect.ValueOf(v); ty.Kind() {
			case reflect.String:
				var valStr string
				vStr := v.(string)
				valStr, ok = val.(string)
				if !ok {
					t.Errorf("%v : %v mismatch type expected %v : found %v", i, k, ty.Kind(), reflect.ValueOf(val).Kind())
				}
				t.Logf("String compare for %v %v=%v", k, vStr, valStr)
				if vStr != valStr {
					t.Errorf("%v : kv value mismatch for %v expected %v : found %v", i, k, vStr, valStr)
				}
			case reflect.Int64:
				var valInt int64
				vInt := v.(int64)
				valInt, ok = val.(int64)
				if !ok {
					t.Errorf("%v : %v mismatch type expected %v : found %v", i, k, ty.Kind(), reflect.ValueOf(val).Kind())
				}
				t.Logf("Int64 compare for %v %v=%v", k, vInt, valInt)
				if vInt != valInt {
					t.Errorf("%v : kv value mismatch for %v expected %v : found %v", i, k, vInt, valInt)
				}
			case reflect.Float64:
				var valFloat float64
				vFloat := v.(float64)
				valFloat, ok = val.(float64)
				if !ok {
					t.Errorf("%v : %v mismatch type expected %v : found %v", i, k, ty.Kind(), reflect.ValueOf(val).Kind())
				}
				t.Logf("Float64 compare for %v %v=%v", k, vFloat, valFloat)
				if vFloat != valFloat {
					t.Errorf("%v : kv value mismatch for %v expected %v : found %v", i, k, vFloat, valFloat)
				}
			default:
				t.Errorf("To be Done for other types : %v", ty.Kind())
			}
		}
	}
}

func Test_readOpens(t *testing.T) {
	p := GrammarParser{}
	samples := [...]string{
		"9\n",                                //0
		"4\n",                                //1
		"4.00000,\nfileSequence436248.m4s\n", //2
	}
	kvs := [...]map[common.AttrId]interface{}{
		{common.INTUnknownAttr: int64(9)},
		{common.INTUnknownAttr: int64(4)},
		{common.INTUnknownAttr: float64(4.00000), common.M3U8Title: "", common.M3U8Uri: "fileSequence436248.m4s"},
	}
	tags := [...]common.TagId{common.M3U8ExtXVersion, common.M3U8TargetDuration, common.M3U8ExtInf}
	remains := [...]int{1, 1, 1}
	states := [...]parserState{searchingTag, searchingTag, searchingTag}
	cols := [...]int{1, 1, 22}
	lines := [...]int{0, 0, 1}
	errors := [...]error{nil, nil, nil,
		errors.New("ERR")}
	for i, sample := range samples {
		if i != 2 {
			continue
		}
		p.state = searchingTag
		p.curTag = tags[i]
		p.line = 0
		p.col = 0
		p.kv = nil
		remain, err := p.readingOpens([]byte(sample))
		if err != nil {
			t.Logf("\n%v: %v", i, err.Error())
		}
		if states[i] != p.state {
			t.Errorf("%v : state expected %v : got %v", i, states[i], p.state)
		}
		if remains[i] != len(remain) {
			t.Errorf("%v : remain expected %v : got %v", i, remains[i], len(remain))
		}
		if (errors[i] != nil && err == nil) || (errors[i] == nil && err != nil) {
			t.Errorf("%v : error expected %v : got %v", i, errors[i], err)
		}
		if cols[i] != p.col {
			t.Errorf("%v : cols expected %v : got %v", i, cols[i], p.col)
		}
		if lines[i] != p.line {
			t.Errorf("%v : lines expected %v : got %v", i, lines[i], p.line)
		}
		for k, v := range kvs[i] {
			var val interface{}
			var ok bool
			if val = p.kv.Get(k); val != nil {
				t.Errorf("%v : kv expected %v : not found", i, k)
			}
			switch ty := reflect.ValueOf(v); ty.Kind() {
			case reflect.String:
				var valStr string
				vStr := v.(string)
				valStr, ok = val.(string)
				if !ok {
					t.Errorf("%v : %v mismatch type expected %v : found %v", i, k, ty.Kind(), reflect.ValueOf(val).Kind())
				}
				t.Logf("String compare for %v %v=%v", k, vStr, valStr)
				if vStr != valStr {
					t.Errorf("%v : kv value mismatch for %v expected %v : found %v", i, k, vStr, valStr)
				}
			case reflect.Int64:
				var valInt int64
				vInt := v.(int64)
				valInt, ok = val.(int64)
				if !ok {
					t.Errorf("%v : %v mismatch type expected %v : found %v", i, k, ty.Kind(), reflect.ValueOf(val).Kind())
				}
				t.Logf("Int64 compare for %v %v=%v", k, vInt, valInt)
				if vInt != valInt {
					t.Errorf("%v : kv value mismatch for %v expected %v : found %v", i, k, vInt, valInt)
				}
			case reflect.Float64:
				var valFloat float64
				vFloat := v.(float64)
				valFloat, ok = val.(float64)
				if !ok {
					t.Errorf("%v : %v mismatch type expected %v : found %v", i, k, ty.Kind(), reflect.ValueOf(val).Kind())
				}
				t.Logf("Float64 compare for %v %v=%v", k, vFloat, valFloat)
				if vFloat != valFloat {
					t.Errorf("%v : kv value mismatch for %v expected %v : found %v", i, k, vFloat, valFloat)
				}
			default:
				t.Errorf("To be Done for other types : %v", ty.Kind())
			}
		}
	}
}
