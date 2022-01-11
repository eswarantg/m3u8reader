package m3u8reader

import (
	"fmt"
	"io"
	"time"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
	"github.com/eswarantg/m3u8reader/parsers/grammarparser"
	"github.com/eswarantg/m3u8reader/parsers/scanparser"
	"github.com/eswarantg/m3u8reader/parsers/yaccparser"
)

type ParserOption int

const (
	M3U8ParserScanner1 ParserOption = iota
	M3U8ParserScanner2
	M3U8ParserScanner3
	M3U8ParserGrammar
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
	buffer                  []byte
}

func (m *M3U8) Done() {
	for _, entry := range m.Entries {
		entry.Done()
	}
}

func (m *M3U8) SetParserOption(opt ParserOption) {
	m.parserOption = opt
}
func (m *M3U8) SetBuffer(buffer []byte) {
	m.buffer = buffer
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
	m.Entries = make([]M3U8Entry, 0, 30)
	m.MediaSequenceNumber = 0
	m.lastSegEntry = nil
	m.lastEntryWCTime = time.Time{}
	m.preloadHintEntry = nil
	m.lastPartWCTime = time.Time{}
}
func (m *M3U8) getParser() parsers.Parser {
	switch m.parserOption {
	case M3U8ParserYacc:
		return &yaccparser.YaccParser{}
	case M3U8ParserScanner2:
		return &scanparser.ScanParser2{}
	case M3U8ParserGrammar:
		return &grammarparser.GrammarParser{}
	case M3U8ParserScanner1:
		fallthrough
	default:
		return &scanparser.ScanParser3{}
	}
}

func (m *M3U8) ParseData(data []byte) (n int, err error) {
	m.Init()
	p := m.getParser()
	if m.buffer != nil {
		p.SetBuffer(m.buffer)
	}
	n, err = p.ParseData(data, m)
	return
}

func (m *M3U8) Read(src io.Reader) (n int, err error) {
	m.Init()
	p := m.getParser()
	if m.buffer != nil {
		p.SetBuffer(m.buffer)
	}
	n, err = p.Parse(src, m)
	return
}

func (m *M3U8) postRecordEntry(entry M3U8Entry) (err error) {
	switch entry.Tag {
	case common.M3U8ExtXPartInf:
		m.partTarget, err = entry.Values.GetFloat64(entry.Tag, common.M3U8PartTarget)
	case common.M3U8ExtXMediaSequence:
		m.MediaSequenceNumber, err = entry.Values.GetInt64(entry.Tag, common.INTUnknownAttr)
		if err != nil {
			return
		}
		m.nextMediaSequenceNumber = m.MediaSequenceNumber
	case common.M3U8TargetDuration:
		m.targetDuration, err = entry.Values.GetInt64(entry.Tag, common.INTUnknownAttr)
	case common.M3U8ExtXIProgramDateTime:
		var t time.Time
		t, err = entry.Values.GetTime(entry.Tag, common.INTUnknownAttr)
		if err != nil {
			return
		}
		m.lastEntryWCTime = t
		m.lastPartWCTime = t
	case common.M3U8ExtInf:
		entry.StoreKV(common.INTProgramDateTime, m.lastEntryWCTime)
		entry.StoreKV(common.INTMediaSequenceNumber, m.nextMediaSequenceNumber)
		m.nextMediaSequenceNumber += 1
		m.nextPartNumber = 0
		var f float64
		f, err = entry.Values.GetFloat64(entry.Tag, common.INTUnknownAttr)
		if err != nil {
			return
		}
		msecDelta := time.Duration(f*1000) * time.Millisecond
		m.lastEntryWCTime = m.lastEntryWCTime.Add(msecDelta)
		m.lastPartWCTime = m.lastEntryWCTime
		m.lastSegEntry = &entry
	case common.M3U8ExtXPart:
		entry.StoreKV(common.INTProgramDateTime, m.lastPartWCTime)
		entry.StoreKV(common.INTMediaSequenceNumber, m.nextMediaSequenceNumber)
		entry.StoreKV(common.INTPartNumber, m.nextPartNumber)
		m.nextPartNumber += 1
		var f float64
		f, err = entry.Values.GetFloat64(entry.Tag, common.M3U8Duration)
		if err != nil {
			return
		}
		msecDelta := time.Duration(f*1000) * time.Millisecond
		m.lastPartWCTime = m.lastPartWCTime.Add(msecDelta)
		m.lastPartEntry = &entry
	case common.M3U8ExtXPreLoadHint:
		//Assuming the lastPartWCTime ith all the XPart data added comuptes to this right start time.
		entry.StoreKV(common.INTProgramDateTime, m.lastPartWCTime)
		m.preloadHintEntry = &entry
	case common.M3U8XSkip:
		//Skip the MediaSequence
		var t int64
		t, err = entry.Values.GetInt64(entry.Tag, common.M3U8SkippedSegments)
		if err != nil {
			return
		}
		m.nextMediaSequenceNumber += t
	}
	m.Entries = append(m.Entries, entry)
	return
}

func (m *M3U8) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
	if kvpairs == nil {
		kvpairs = parsers.AttrKVPairsPool.Get().(*parsers.AttrKVPairs)
	}
	entry := M3U8Entry{Tag: tag, Values: kvpairs}
	return m.postRecordEntry(entry)
}

func (m *M3U8) GetVideoMediaPlaylist(maxBitRateBps int64) (toret *M3U8Entry, err error) {
	toret = nil
	var entryObj M3U8Entry
	curSelectBW := int64(-1)
	for _, entry := range m.Entries {
		if entry.Tag == common.M3U8ExtXStreamInf {
			entryBW := entry.Values.Get(common.M3U8Bandwidth).(int64)
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
		if entry.Tag == common.M3U8ExtXMedia {
			if lang == entry.Values.Get(common.M3U8Language).(string) {
				if vidEntry.Values.Get(common.M3U8Audio).(string) == entry.Values.Get(common.M3U8GroupId).(string) {
					toret = &entry
					break
				}
			}
		}
	}
	return
}
