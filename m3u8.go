package m3u8reader

import (
	"fmt"
	"io"
	"time"
)

type M3U8Entry struct {
	Tag    TagId
	Values map[AttrId]interface{}
}

func (m *M3U8Entry) storeKV(k AttrId, v interface{}) {
	if m.Values == nil {
		m.Values = make(map[AttrId]interface{})
	}
	m.Values[k] = v
}

func (m *M3U8Entry) String() string {
	return fmt.Sprintf("%v %v", m.Tag, m.Values)
}

func (m *M3U8Entry) URI() (string, error) {
	switch m.Tag {
	case M3U8ExtXStreamInf:
		return m.Values[INTUnknownAttr].(string), nil
	case M3U8ExtXMedia:
		return m.Values[M3U8Uri].(string), nil
	case M3U8ExtInf:
		return m.Values[M3U8Uri].(string), nil
	case M3U8ExtXPreLoadHint:
		return m.Values[M3U8Uri].(string), nil
	}
	return "", fmt.Errorf("URI not available")
}

type ParserOption int

const (
	M3U8ParserQuotesSafe ParserOption = iota
	M3U8ParserQuotesUnsafe
	M3U8ParserYacc
)

type M3U8 struct {
	Entries                 []M3U8Entry
	MediaSequenceNumber     int64
	targetDuration          int64
	partTarget              float64
	lastSegEntry            *M3U8Entry
	lastPartEntry           *M3U8Entry
	lastEntryWCTime         time.Time
	lastPartWCTime          time.Time
	preloadHintEntry        *M3U8Entry
	nextMediaSequenceNumber int64
	nextPartNumber          int64
	parserOption            ParserOption
}

func (m *M3U8) SetParserOption(opt ParserOption) {
	m.parserOption = opt
}

func (m *M3U8) String() string {
	toret := ""
	for _, entry := range m.Entries {
		toret += fmt.Sprintf("\n%v", entry.String())
	}
	return toret
}

func (m *M3U8) TargetDuration() int64 {
	return m.targetDuration
}
func (m *M3U8) PartTarget() float64 {
	return m.partTarget
}

func (m *M3U8) LastSegment() *M3U8Entry {
	return m.lastSegEntry
}
func (m *M3U8) LastSegmentTime() time.Time {
	return m.lastEntryWCTime
}
func (m *M3U8) PreloadHintEntry() *M3U8Entry {
	return m.preloadHintEntry
}
func (m *M3U8) LastPart() *M3U8Entry {
	return m.lastPartEntry
}

func (m *M3U8) LastPartTime() time.Time {
	return m.lastPartWCTime
}

func (m *M3U8) Init() {
	m.Entries = make([]M3U8Entry, 0)
	m.MediaSequenceNumber = 0
	m.lastSegEntry = nil
	m.lastEntryWCTime = time.Time{}
	m.preloadHintEntry = nil
	m.lastPartWCTime = time.Time{}
}

func (m *M3U8) yyParse(src io.Reader) (err error) {
	yyErrorVerbose = true
	lex := NewLexerWithInit(src, func(l *Lexer) {
		l.parseResult = m
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
	if result != 0 {
		err = fmt.Errorf("yyparse returned failure")
	}
	return
}

func (m *M3U8) Read(src io.Reader) (n int, err error) {
	m.Init()
	switch m.parserOption {
	case M3U8ParserYacc:
		err = m.yyParse(src)
	case M3U8ParserQuotesUnsafe:
		n, err = parseM3U8_fast(src, m)
	case M3U8ParserQuotesSafe:
		fallthrough
	default:
		n, err = parseM3U8(src, m)
	}
	return
}

func (m *M3U8) postRecordEntry(entry M3U8Entry) (err error) {
	m.Entries = append(m.Entries, entry)
	switch entry.Tag {
	case M3U8ExtXPartInf:
		m.partTarget = entry.Values[M3U8PartTarget].(float64)
	case M3U8ExtXMediaSequence:
		m.MediaSequenceNumber = entry.Values[INTUnknownAttr].(int64)
		m.nextMediaSequenceNumber = m.MediaSequenceNumber
	case M3U8TargetDuration:
		m.targetDuration = entry.Values[INTUnknownAttr].(int64)
	case M3U8ExtXIProgramDateTime:
		m.lastEntryWCTime = entry.Values[INTUnknownAttr].(time.Time)
		m.lastPartWCTime = entry.Values[INTUnknownAttr].(time.Time)
	case M3U8ExtInf:
		entry.Values[INTProgramDateTime] = m.lastEntryWCTime
		entry.Values[INTMediaSequenceNumber] = m.nextMediaSequenceNumber
		m.nextMediaSequenceNumber += 1
		m.nextPartNumber = 0
		msecDelta := time.Duration(entry.Values[INTUnknownAttr].(float64)*1000) * time.Millisecond
		m.lastEntryWCTime = m.lastEntryWCTime.Add(msecDelta)
		m.lastPartWCTime = m.lastEntryWCTime
		m.lastSegEntry = &entry
	case M3U8ExtXPart:
		entry.Values[INTProgramDateTime] = m.lastPartWCTime
		entry.Values[INTMediaSequenceNumber] = m.nextMediaSequenceNumber
		entry.Values[INTPartNumber] = m.nextPartNumber
		m.nextPartNumber += 1
		msecDelta := time.Duration(entry.Values[M3U8Duration].(float64)*1000) * time.Millisecond
		m.lastPartWCTime = m.lastPartWCTime.Add(msecDelta)
		m.lastPartEntry = &entry
	case M3U8ExtXPreLoadHint:
		//Assuming the lastPartWCTime ith all the XPart data added comuptes to this right start time.
		entry.Values[INTProgramDateTime] = m.lastPartWCTime
		m.preloadHintEntry = &entry
	case M3U8XSkip:
		//Skip the MediaSequence
		m.nextMediaSequenceNumber += entry.Values[M3U8SkippedSegments].(int64)
	}
	return
}

func (m *M3U8) postRecord(tag TagId, kvpairs map[AttrId]interface{}) (err error) {
	entry := M3U8Entry{Tag: tag, Values: kvpairs}
	err = m.decorateEntry(&entry)
	if err != nil {
		return
	}
	return m.postRecordEntry(entry)
}

func (m *M3U8) decorateEntry(entry *M3U8Entry) (err error) {
	if decorateFn, ok := decorators[entry.Tag]; ok {
		err = decorateFn(entry)
	}
	return
}

func (m *M3U8) GetVideoMediaPlaylist(maxBitRateBps int64) (toret *M3U8Entry, err error) {
	toret = nil
	var entryObj M3U8Entry
	curSelectBW := int64(-1)
	for _, entry := range m.Entries {
		if entry.Tag == M3U8ExtXStreamInf {
			entryBW := entry.Values[M3U8Bandwidth].(int64)
			if entryBW <= maxBitRateBps && entryBW > curSelectBW {
				entryObj = entry
				toret = &entryObj
				curSelectBW = entryBW
			}
		}
	}
	return toret, err
}

func (m *M3U8) GetAudioMediaPlaylist(vidEntry M3U8Entry, lang string) (toret *M3U8Entry, err error) {
	toret = nil
	for _, entry := range m.Entries {
		if entry.Tag == M3U8ExtXMedia {
			if lang == entry.Values[M3U8Language].(string) {
				if vidEntry.Values[M3U8Audio].(string) == entry.Values[M3U8GroupId].(string) {
					toret = &entry
					break
				}
			}
		}
	}
	return
}
