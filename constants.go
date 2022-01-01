package m3u8reader

var tagNames = [...]string{
	"EXTM3U",
	"EXT-X-VERSION",
	"EXT-X-INDEPENDENT-SEGMENTS",
	"EXT-X-MEDIA",
	"EXT-X-STREAM-INF",
	"EXT-X-TARGETDURATION",
	"EXT-X-SERVER-CONTROL",
	"EXT-X-PART-INF",
	"EXT-X-MEDIA-SEQUENCE",
	"EXT-X-SKIP",
	"EXTINF",
	"EXT-X-PROGRAM-DATE-TIME",
	"EXT-X-PART",
	"EXT-X-PRELOAD-HINT",
	"EXT-X-RENDITION-REPORT",
	"EXT-X-MAP",
}

var attrNames = [...]string{
	"BANDWIDTH",
	"AVERAGE-BANDWIDTH",
	"RESOLUTION",
	"FRAME-RATE",
	"CODECS",
	"AUDIO",
	"TYPE",
	"GROUP-ID",
	"NAME",
	"DEFAULT",
	"AUTOSELECT",
	"LANGUAGE",
	"CHANNELS",
	"URI",
	"CAN-BLOCK-RELOAD",
	"CAN-SKIP-UNTIL",
	"PART-HOLD-BACK",
	"PART-TARGET",
	"SKIPPED-SEGMENTS",
	"DURATION",
	"INDEPENDENT",
	"LAST-MSN",
	"LAST-PART",
}

func tagName(token int) string {
	token -= TAG_FIRST + 1
	return tagNames[token]
}
func attrName(token int) string {
	token -= ATTR_FIRST + 1
	return attrNames[token]
}

const (
	M3U8FormatIdentifier        = "EXTM3U"
	M3U8ExtXVersion             = "EXT-X-VERSION"
	M3U8TargetDuration          = "EXT-X-TARGETDURATION"
	M3U8ExtXMedia               = "EXT-X-MEDIA"
	M3U8ExtXStreamInf           = "EXT-X-STREAM-INF"
	M3U8ExtXIFrameStreamInf     = "EXT-X-I-FRAME-STREAM-INF"
	M3U8ExtXMediaSequence       = "EXT-X-MEDIA-SEQUENCE"
	M3U8ExtXIProgramDateTime    = "EXT-X-PROGRAM-DATE-TIME"
	M3U8ExtInf                  = "EXTINF"
	M3U8ExtXIndependentSegments = "EXT-X-INDEPENDENT-SEGMENTS"
	M3U8ExtXPart                = "EXT-X-PART"
	M3U8ExtXPartInf             = "EXT-X-PART-INF"
	M3U8ExtXRenditionReport     = "EXT-X-RENDITION-REPORT"
	M3U8ExtXPreLoadHint         = "EXT-X-PRELOAD-HINT"
	M3U8ExtXServerControl       = "EXT-X-SERVER-CONTROL"
	M3U8ExtXDiscontinuity       = "EXT-X-DISCONTINUITY"
	M3U8ExtXEndList             = "EXT-X-ENDLIST"
	M3U8ExtXPlaylistType        = "EXT-X-PLAYLIST-TYPE"
	M3U8XSkip                   = "EXT-X-SKIP"
)
