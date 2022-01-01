package m3u8reader

import (
	"fmt"
	"io"
	"time"
)

type M3U8Entry struct {
	Tag    string
	Values map[string]interface{}
}

func (m *M3U8Entry) StoreKV(k string, v interface{}) {
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[k] = v
}

func (m *M3U8Entry) String() string {
	return fmt.Sprintf("%v %v", m.Tag, m.Values)
}

func (m *M3U8Entry) URI() (string, error) {
	switch m.Tag {
	case M3U8ExtXStreamInf:
		return m.Values[m3u8UnknownKey].(string), nil
	case M3U8ExtXMedia:
		return m.Values["URI"].(string), nil
	case M3U8ExtInf:
		return m.Values["URI"].(string), nil
	case M3U8ExtXPreLoadHint:
		return m.Values["URI"].(string), nil
	}
	return "", fmt.Errorf("URI not available")
}

type ParserOption int

const (
	M3U8ParserQuotesSafe ParserOption = iota
	M3U8ParserQuotesUnsafe
)

type M3U8 struct {
	Entries                 []M3U8Entry
	MediaSequenceNumber     int64
	targetDuration          float64
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

func (m *M3U8) TargetDuration() float64 {
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

func (m *M3U8) Read(src io.Reader) (n int, err error) {
	m.Entries = make([]M3U8Entry, 0)
	m.MediaSequenceNumber = 0
	m.lastSegEntry = nil
	m.lastEntryWCTime = time.Time{}
	m.preloadHintEntry = nil
	m.lastPartWCTime = time.Time{}
	switch m.parserOption {
	case M3U8ParserQuotesUnsafe:
		n, err = parseM3U8_fast(src, m)
	case M3U8ParserQuotesSafe:
		fallthrough
	default:
		n, err = parseM3U8(src, m)
	}
	return
}

func (m *M3U8) PostRecordEntry(entry M3U8Entry) (err error) {
	m.Entries = append(m.Entries, entry)
	switch entry.Tag {
	case M3U8ExtXPartInf:
		m.partTarget = entry.Values["PART-TARGET"].(float64)
	case M3U8ExtXMediaSequence:
		m.MediaSequenceNumber = entry.Values[m3u8UnknownKey].(int64)
		m.nextMediaSequenceNumber = m.MediaSequenceNumber
	case M3U8TargetDuration:
		m.targetDuration = entry.Values[m3u8UnknownKey].(float64)
	case M3U8ExtXIProgramDateTime:
		m.lastEntryWCTime = entry.Values[m3u8UnknownKey].(time.Time)
		m.lastPartWCTime = entry.Values[m3u8UnknownKey].(time.Time)
	case M3U8ExtInf:
		entry.Values["programDateTime"] = m.lastEntryWCTime
		entry.Values["mediaSequenceNumber"] = m.nextMediaSequenceNumber
		m.nextMediaSequenceNumber += 1
		m.nextPartNumber = 0
		msecDelta := time.Duration(entry.Values[m3u8UnknownKey].(float64)*1000) * time.Millisecond
		m.lastEntryWCTime = m.lastEntryWCTime.Add(msecDelta)
		m.lastPartWCTime = m.lastEntryWCTime
		m.lastSegEntry = &entry
	case M3U8ExtXPart:
		entry.Values["programDateTime"] = m.lastPartWCTime
		entry.Values["mediaSequenceNumber"] = m.nextMediaSequenceNumber
		entry.Values["partNumber"] = m.nextPartNumber
		m.nextPartNumber += 1
		msecDelta := time.Duration(entry.Values["DURATION"].(float64)*1000) * time.Millisecond
		m.lastPartWCTime = m.lastPartWCTime.Add(msecDelta)
		m.lastPartEntry = &entry
	case M3U8ExtXPreLoadHint:
		//Assuming the lastPartWCTime ith all the XPart data added comuptes to this right start time.
		entry.Values["programDateTime"] = m.lastPartWCTime
		m.preloadHintEntry = &entry
	case M3U8XSkip:
		//Skip the MediaSequence
		m.nextMediaSequenceNumber += entry.Values["SKIPPED-SEGMENTS"].(int64)
	}
	return
}

func (m *M3U8) postRecord(tag string, kvpairs map[string]interface{}) (err error) {
	entry := M3U8Entry{Tag: tag, Values: kvpairs}
	err = m.decorateEntry(&entry)
	if err != nil {
		return
	}
	return m.PostRecordEntry(entry)
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
			entryBW := entry.Values["BANDWIDTH"].(int64)
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
			if lang == entry.Values["LANGUAGE"].(string) {
				if vidEntry.Values["AUDIO"].(string) == entry.Values["GROUP-ID"].(string) {
					toret = &entry
					break
				}
			}
		}
	}
	return
}
