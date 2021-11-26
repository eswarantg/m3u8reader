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
	}
	return "", fmt.Errorf("URI not available")
}

type M3U8 struct {
	Entries             []M3U8Entry
	MediaSequenceNumber int64
	lastEntry           *M3U8Entry
	lastEntryWCTime     time.Time
}

func (m *M3U8) String() string {
	toret := ""
	for _, entry := range m.Entries {
		toret += fmt.Sprintf("\n%v", entry.String())
	}
	return toret
}

func (m *M3U8) LastSegment() *M3U8Entry {
	return m.lastEntry
}
func (m *M3U8) LastSegmentTime() time.Time {
	return m.lastEntryWCTime
}

func (m *M3U8) Read(src io.Reader) (n int, err error) {
	m.Entries = make([]M3U8Entry, 0)
	return parseM3U8(src, m)
}

func (m *M3U8) postRecord(tag string, kvpairs map[string]interface{}) (err error) {
	entry := M3U8Entry{Tag: tag, Values: kvpairs}
	err = m.decorateEntry(&entry)
	if err != nil {
		return
	}
	m.Entries = append(m.Entries, entry)
	switch entry.Tag {
	case M3U8ExtXMediaSequence:
		m.MediaSequenceNumber = entry.Values[m3u8UnknownKey].(int64)
	case M3U8ExtXIProgramDateTime:
		m.lastEntryWCTime = entry.Values[m3u8UnknownKey].(time.Time)
	case M3U8ExtInf:
		m.lastEntry = &entry
	}
	return
}

func (m *M3U8) decorateEntry(entry *M3U8Entry) (err error) {
	if decorateFn, ok := decorators[entry.Tag]; ok {
		err = decorateFn(entry)
	}
	return
}

func (m *M3U8) GetVideoMediaPlaylist(maxBitRateBps int64) (toret *M3U8Entry, err error) {
	toret = nil
	curSelectBW := int64(-1)
	for _, entry := range m.Entries {
		if entry.Tag == M3U8ExtXStreamInf {
			entryBW := entry.Values["BANDWIDTH"].(int64)
			if entryBW <= maxBitRateBps && entryBW > curSelectBW {
				toret = &entry
				curSelectBW = entryBW
			}
		}
	}
	return
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
