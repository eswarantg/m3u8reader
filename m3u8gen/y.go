// Code generated by goyacc m3u8.y. DO NOT EDIT.

//line m3u8.y:2
package m3u8gen

import __yyfmt__ "fmt"

//line m3u8.y:2

import "github.com/eswarantg/m3u8reader"
import "time"

func TagName(token int) string {
	const TOKSTART = 4
	token -= TAG_FIRST + 1
	token += TOKSTART
	return yyToknames[token]
}
func AttrName(token int) string {
	const TOKSTART = 4
	token -= ATTR_FIRST + 1
	token += TOKSTART
	return yyToknames[token]
}

func setResult(l yyLexer, v *m3u8reader.M3U8) {
	l.(*Lexer).parseResult = v
}

//line m3u8.y:26
type yySymType struct {
	yys      int
	i        int
	i64      int64
	f        float64
	s        string
	r        string
	t        time.Time
	val      interface{}
	entry    m3u8reader.M3U8Entry
	manifest m3u8reader.M3U8
}

const TAG_FIRST = 57346
const TAG_EXTM3U = 57347
const TAG_EXT_X_VERSION = 57348
const TAG_EXT_X_INDEPENDENT_SEGMENTS = 57349
const TAG_EXT_X_MEDIA = 57350
const TAG_EXT_X_STREAM_INF = 57351
const TAG_EXT_X_TARGETDURATION = 57352
const TAG_EXT_X_SERVER_CONTROL = 57353
const TAG_EXT_X_PART_INF = 57354
const TAG_EXT_X_MEDIA_SEQUENCE = 57355
const TAG_EXT_X_SKIP = 57356
const TAG_EXTINF = 57357
const TAG_EXT_X_PROGRAM_DATE_TIME = 57358
const TAG_EXT_X_PART = 57359
const TAG_EXT_X_PRELOAD_HINT = 57360
const TAG_EXT_X_RENDITION_REPORT = 57361
const TAG_EXT_X_MAP = 57362
const COMMA = 57363
const EQUALTO = 57364
const SECONDLINEVALUE = 57365
const ATTR_FIRST = 57366
const ATTR_BANDWIDTH = 57367
const ATTR_AVERAGE_BANDWIDTH = 57368
const ATTR_RESOLUTION = 57369
const ATTR_FRAME_RATE = 57370
const ATTR_CODECS = 57371
const ATTR_AUDIO = 57372
const ATTR_TYPE = 57373
const ATTR_GROUP_ID = 57374
const ATTR_NAME = 57375
const ATTR_DEFAULT = 57376
const ATTR_AUTOSELECT = 57377
const ATTR_LANGUAGE = 57378
const ATTR_CHANNELS = 57379
const ATTR_URI = 57380
const ATTR_CAN_BLOCK_RELOAD = 57381
const ATTR_CAN_SKIP_UNTIL = 57382
const ATTR_PART_HOLD_BACK = 57383
const ATTR_PART_TARGET = 57384
const ATTR_SKIPPED_SEGMENTS = 57385
const ATTR_DURATION = 57386
const ATTR_INDEPENDENT = 57387
const ATTR_LAST_MSN = 57388
const ATTR_LAST_PART = 57389
const ATTRKEY = 57390
const INTEGERVAL = 57391
const FLOATVAL = 57392
const STRINGVAL = 57393
const RESOLUTIONVAL = 57394
const TIMEVAL = 57395

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TAG_FIRST",
	"TAG_EXTM3U",
	"TAG_EXT_X_VERSION",
	"TAG_EXT_X_INDEPENDENT_SEGMENTS",
	"TAG_EXT_X_MEDIA",
	"TAG_EXT_X_STREAM_INF",
	"TAG_EXT_X_TARGETDURATION",
	"TAG_EXT_X_SERVER_CONTROL",
	"TAG_EXT_X_PART_INF",
	"TAG_EXT_X_MEDIA_SEQUENCE",
	"TAG_EXT_X_SKIP",
	"TAG_EXTINF",
	"TAG_EXT_X_PROGRAM_DATE_TIME",
	"TAG_EXT_X_PART",
	"TAG_EXT_X_PRELOAD_HINT",
	"TAG_EXT_X_RENDITION_REPORT",
	"TAG_EXT_X_MAP",
	"COMMA",
	"EQUALTO",
	"SECONDLINEVALUE",
	"ATTR_FIRST",
	"ATTR_BANDWIDTH",
	"ATTR_AVERAGE_BANDWIDTH",
	"ATTR_RESOLUTION",
	"ATTR_FRAME_RATE",
	"ATTR_CODECS",
	"ATTR_AUDIO",
	"ATTR_TYPE",
	"ATTR_GROUP_ID",
	"ATTR_NAME",
	"ATTR_DEFAULT",
	"ATTR_AUTOSELECT",
	"ATTR_LANGUAGE",
	"ATTR_CHANNELS",
	"ATTR_URI",
	"ATTR_CAN_BLOCK_RELOAD",
	"ATTR_CAN_SKIP_UNTIL",
	"ATTR_PART_HOLD_BACK",
	"ATTR_PART_TARGET",
	"ATTR_SKIPPED_SEGMENTS",
	"ATTR_DURATION",
	"ATTR_INDEPENDENT",
	"ATTR_LAST_MSN",
	"ATTR_LAST_PART",
	"ATTRKEY",
	"INTEGERVAL",
	"FLOATVAL",
	"STRINGVAL",
	"RESOLUTIONVAL",
	"TIMEVAL",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line m3u8.y:162

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 89

var yyAct = [...]int{
	26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45,
	46, 47, 48, 25, 23, 65, 66, 67, 68, 69,
	57, 56, 55, 53, 50, 21, 5, 7, 8, 6,
	9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
	19, 22, 64, 63, 75, 62, 74, 63, 72, 71,
	49, 2, 51, 52, 1, 54, 3, 24, 58, 59,
	60, 61, 4, 0, 0, 0, 20, 0, 70, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 73,
}

var yyPact = [...]int{
	56, -1000, 30, 30, -1000, -14, -25, -1000, -25, -15,
	-25, -25, -16, -25, -18, -23, -25, -25, -25, -25,
	-1000, -1000, 32, -1000, -24, -24, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 36,
	-1000, 36, 36, -1000, 36, 38, 37, -1000, 36, 36,
	36, 36, -1000, -25, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 33, 31, -1000, -1000, -1000,
}

var yyPgo = [...]int{
	0, 67, 52, 51, 24, 72, 66, 64,
}

var yyR1 = [...]int{
	0, 7, 6, 6, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	3, 3, 4, 4, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
	2, 2,
}

var yyR2 = [...]int{
	0, 2, 1, 2, 2, 3, 1, 2, 2, 2,
	2, 2, 2, 4, 4, 2, 2, 2, 2, 2,
	1, 3, 2, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1,
}

var yyChk = [...]int{
	-1000, -7, 5, -6, -5, 6, 9, 7, 8, 10,
	11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	-5, 49, -3, -4, -1, 48, 25, 26, 27, 28,
	29, 30, 31, 32, 33, 34, 35, 36, 37, 38,
	39, 40, 41, 42, 43, 44, 45, 46, 47, -3,
	49, -3, -3, 49, -3, 50, 49, 53, -3, -3,
	-3, -3, 23, 21, -2, 49, 50, 51, 52, 53,
	-2, 21, 21, -4, 23, 23,
}

var yyDef = [...]int{
	0, -2, 0, 1, 2, 0, 0, 6, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	3, 4, 0, 20, 0, 0, 24, 25, 26, 27,
	28, 29, 30, 31, 32, 33, 34, 35, 36, 37,
	38, 39, 40, 41, 42, 43, 44, 45, 46, 7,
	8, 9, 10, 11, 12, 0, 0, 15, 16, 17,
	18, 19, 5, 0, 22, 47, 48, 49, 50, 51,
	23, 0, 0, 21, 13, 14,
}

var yyTok1 = [...]int{
	1,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53,
}

var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:105
		{
			setResult(yylex, &yyVAL.manifest)
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:107
		{
			yyVAL.manifest.PostRecordEntry(yyDollar[1].entry)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:108
		{
			yyVAL.manifest.PostRecordEntry(yyDollar[2].entry)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:110
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", yyDollar[2].i64)
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
//line m3u8.y:111
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", yyDollar[3].s)
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:112
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:113
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:114
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", yyDollar[2].i64)
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:115
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:116
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:117
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", yyDollar[2].i64)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:118
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
//line m3u8.y:119
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", yyDollar[2].f)
			yyVAL.entry.StoreKV("URI", yyDollar[4].s)
		}
	case 14:
		yyDollar = yyS[yypt-4 : yypt+1]
//line m3u8.y:120
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", float64(yyDollar[2].i64))
			yyVAL.entry.StoreKV("URI", yyDollar[4].s)
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:121
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
			yyVAL.entry.StoreKV("#", yyDollar[2].t)
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:122
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:123
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:124
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:125
		{
			yyVAL.entry.Tag = TagName(yyDollar[1].i)
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:130
		{
			yyVAL.entry.StoreKV(AttrName(yyDollar[1].i), yyDollar[2].val)
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
//line m3u8.y:131
		{
			yyVAL.entry.StoreKV(yyDollar[1].s, yyDollar[2].val)
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:133
		{
			yyVAL.i = yyDollar[1].i
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:134
		{
			yyVAL.i = yyDollar[1].i
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:135
		{
			yyVAL.i = yyDollar[1].i
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:136
		{
			yyVAL.i = yyDollar[1].i
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:137
		{
			yyVAL.i = yyDollar[1].i
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:138
		{
			yyVAL.i = yyDollar[1].i
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:139
		{
			yyVAL.i = yyDollar[1].i
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:140
		{
			yyVAL.i = yyDollar[1].i
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:141
		{
			yyVAL.i = yyDollar[1].i
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:142
		{
			yyVAL.i = yyDollar[1].i
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:143
		{
			yyVAL.i = yyDollar[1].i
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:144
		{
			yyVAL.i = yyDollar[1].i
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:145
		{
			yyVAL.i = yyDollar[1].i
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:146
		{
			yyVAL.i = yyDollar[1].i
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:147
		{
			yyVAL.i = yyDollar[1].i
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:148
		{
			yyVAL.i = yyDollar[1].i
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:149
		{
			yyVAL.i = yyDollar[1].i
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:150
		{
			yyVAL.i = yyDollar[1].i
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:151
		{
			yyVAL.i = yyDollar[1].i
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:152
		{
			yyVAL.i = yyDollar[1].i
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:153
		{
			yyVAL.i = yyDollar[1].i
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:154
		{
			yyVAL.i = yyDollar[1].i
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:155
		{
			yyVAL.i = yyDollar[1].i
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:157
		{
			yyVAL.val = yyDollar[1].i64
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:158
		{
			yyVAL.val = yyDollar[1].f
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:159
		{
			yyVAL.val = yyDollar[1].s
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:160
		{
			yyVAL.val = yyDollar[1].r
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
//line m3u8.y:161
		{
			yyVAL.val = yyDollar[1].t
		}
	}
	goto yystack /* stack new state and value */
}
