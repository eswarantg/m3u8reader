package m3u8reader

//Types of Data/Value to be parsed
type ValueType int

const (
	valueOptional                 ValueType = 1
	valueDecimalInt                         = 1 << 1
	valueHexaDecimalSeq                     = 1 << 2
	valueUnSignedDecimalFloat               = 1 << 3
	valueSignedDecimalFloat                 = 1 << 4
	valueQuotedString                       = 1 << 5
	valueEnumeratedString                   = 1 << 6
	valueDecimalResolution                  = 1 << 7
	valueNextLineEnumeratedString           = 1 << 8
	valueDateTime                           = 1 << 9
	valueUTF8Text                           = 1 << 10
)

//All TagNames in M3U8
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
	"EXT-X-I-FRAME-STREAM-INF",
	"EXT-X-DISCONTINUITY",
	"EXT-X-ENDLIST",
	"EXT-X-PLAYLIST-TYPE",
	"EXT-X-BYTERANGE",
	"EXT-X-KEY",
	"EXT-X-DATERANGE",
	"EXT-X-DISCONTINUITY-SEQUENCE",
	"EXT-X-I-FRAMES-ONLY",
	"EXT-X-SESSION-KEY",
	"EXT-X-SESSION-DATA",
	"EXT-X-START",
}

//Internal Identification Number of each Tag
//To avoid storing/comparing tag
type TagId int

const (
	M3U8FormatIdentifier TagId = iota
	M3U8ExtXVersion
	M3U8ExtXIndependentSegments
	M3U8ExtXMedia
	M3U8ExtXStreamInf
	M3U8TargetDuration
	M3U8ExtXServerControl
	M3U8ExtXPartInf
	M3U8ExtXMediaSequence
	M3U8XSkip
	M3U8ExtInf
	M3U8ExtXIProgramDateTime
	M3U8ExtXPart
	M3U8ExtXPreLoadHint
	M3U8ExtXRenditionReport
	M3U8ExtXMap
	M3U8ExtXIFrameStreamInf
	M3U8ExtXDiscontinuity
	M3U8ExtXEndList
	M3U8ExtXPlaylistType
	M3U8ExtXByteRange
	M3U8ExtXKey
	M3U8ExtXDataRange
	M3U8ExtXDiscontinuitySequence
	M3U8ExtXIFramesOnly
	M3U8ExtXSesionKey
	M3U8ExtXSessionData
	M3U8ExtXStart
)

var tagToTagId map[string]TagId = map[string]TagId{
	"EXTM3U":                       M3U8FormatIdentifier,
	"EXT-X-VERSION":                M3U8ExtXVersion,
	"EXT-X-INDEPENDENT-SEGMENTS":   M3U8ExtXIndependentSegments,
	"EXT-X-MEDIA":                  M3U8ExtXMedia,
	"EXT-X-STREAM-INF":             M3U8ExtXStreamInf,
	"EXT-X-TARGETDURATION":         M3U8TargetDuration,
	"EXT-X-SERVER-CONTROL":         M3U8ExtXServerControl,
	"EXT-X-PART-INF":               M3U8ExtXPartInf,
	"EXT-X-MEDIA-SEQUENCE":         M3U8ExtXMediaSequence,
	"EXT-X-SKIP":                   M3U8XSkip,
	"EXTINF":                       M3U8ExtInf,
	"EXT-X-PROGRAM-DATE-TIME":      M3U8ExtXIProgramDateTime,
	"EXT-X-PART":                   M3U8ExtXPart,
	"EXT-X-PRELOAD-HINT":           M3U8ExtXPreLoadHint,
	"EXT-X-RENDITION-REPORT":       M3U8ExtXRenditionReport,
	"EXT-X-MAP":                    M3U8ExtXMap,
	"EXT-X-I-FRAME-STREAM-INF":     M3U8ExtXIFrameStreamInf,
	"EXT-X-DISCONTINUITY":          M3U8ExtXDiscontinuity,
	"EXT-X-ENDLIST":                M3U8ExtXEndList,
	"EXT-X-PLAYLIST-TYPE":          M3U8ExtXPlaylistType,
	"EXT-X-BYTERANGE":              M3U8ExtXByteRange,
	"EXT-X-KEY":                    M3U8ExtXKey,
	"EXT-X-DATERANGE":              M3U8ExtXDataRange,
	"EXT-X-DISCONTINUITY-SEQUENCE": M3U8ExtXDiscontinuitySequence,
	"EXT-X-I-FRAMES-ONLY":          M3U8ExtXIFramesOnly,
	"EXT-X-SESSION-KEY":            M3U8ExtXSesionKey,
	"EXT-X-SESSION-DATA":           M3U8ExtXSessionData,
	"EXT-X-START":                  M3U8ExtXStart,
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
	"METHOD",
	"IV",
	"KEYFORMAT",
	"KEYFORMATVERSIONS",
	"BYTERANGE",
	"ID",
	"CLASS",
	"START-DATE",
	"END-DATE",
	"PLANNED-DURATION",
	"SCTE35-CMD",
	"SCTE35-OUT",
	"SCTE35-IN",
	"END-ON-NEXT",
	"ASSOC-LANGUAGE",
	"FORCED",
	"INSTREAM-ID",
	"CHARACTERISTICS",
	"HDCP-LEVEL",
	"VIDEO",
	"SUBTITLES",
	"CLOSED-CAPTIONS",
	"TIME-OFFSET",
	"PRECISE",
	"DATA-ID",
	"VALUE",
	"#",
	"programDataTime",
	"mediaSequenceNumber",
	"partNumber",
}

//To avoid storing/comparing Attr
type AttrId int

const (
	M3U8Bandwidth AttrId = iota
	M3U8AverageBandwidth
	M3U8Resolution
	M3U8FrameRate
	M3U8Codecs
	M3U8Audio
	M3U8Type
	M3U8GroupId
	M3U8Name
	M3U8Default
	M3U8AutoSelect
	M3U8Language
	M3U8Channels
	M3U8Uri
	M3U8CanBlockReload
	M3U8CanSkipUntil
	M3U8PartHoldBack
	M3U8PartTarget
	M3U8SkippedSegments
	M3U8Duration
	M3U8Independent
	M3U8LastMsn
	M3U8LastPart
	M3U8Method
	M3U8IV
	M3U8KeyFormat
	M3U8KeyFormatVersions
	M3U8ByteRange
	M3U8Id
	M3U8Class
	M3U8StartDate
	M3U8EndDate
	M3U8PlannedDuration
	M3U8Scte35Cmd
	M3U8Scte35Out
	M3U8Scte35In
	M3U8EndOnNext
	M3U8AssocLanguage
	M3U8Forced
	M3U8InStreamId
	M3U8Characteristics
	M3U8HdcpLevel
	M3U8Video
	M3U8Subtitles
	M3U8ClosedCaptions
	M3U8TimeOffset
	M3U8Precise
	M3U8DataId
	M3U8Value

	INTUnknownAttr
	INTProgramDateTime
	INTMediaSequenceNumber
	INTPartNumber
)

var attrToAttrId map[string]AttrId = map[string]AttrId{
	"BANDWIDTH":           M3U8Bandwidth,
	"AVERAGE-BANDWIDTH":   M3U8AverageBandwidth,
	"RESOLUTION":          M3U8Resolution,
	"FRAME-RATE":          M3U8FrameRate,
	"CODECS":              M3U8Codecs,
	"AUDIO":               M3U8Audio,
	"TYPE":                M3U8Type,
	"GROUP-ID":            M3U8GroupId,
	"NAME":                M3U8Name,
	"DEFAULT":             M3U8Default,
	"AUTOSELECT":          M3U8AutoSelect,
	"LANGUAGE":            M3U8Language,
	"CHANNELS":            M3U8Channels,
	"URI":                 M3U8Uri,
	"CAN-BLOCK-RELOAD":    M3U8CanBlockReload,
	"CAN-SKIP-UNTIL":      M3U8CanSkipUntil,
	"PART-HOLD-BACK":      M3U8PartHoldBack,
	"PART-TARGET":         M3U8PartTarget,
	"SKIPPED-SEGMENTS":    M3U8SkippedSegments,
	"DURATION":            M3U8Duration,
	"INDEPENDENT":         M3U8Independent,
	"LAST-MSN":            M3U8LastMsn,
	"LAST-PART":           M3U8LastPart,
	"METHOD":              M3U8Method,
	"IV":                  M3U8IV,
	"KEYFORMAT":           M3U8KeyFormat,
	"KEYFORMATVERSIONS":   M3U8KeyFormatVersions,
	"BYTERANGE":           M3U8ByteRange,
	"ID":                  M3U8Id,
	"CLASS":               M3U8Class,
	"START-DATE":          M3U8StartDate,
	"END-DATE":            M3U8EndDate,
	"PLANNED-DURATION":    M3U8PlannedDuration,
	"SCTE35-CMD":          M3U8Scte35Cmd,
	"SCTE35-OUT":          M3U8Scte35Out,
	"SCTE35-IN":           M3U8Scte35In,
	"END-ON-NEXT":         M3U8EndOnNext,
	"ASSOC-LANGUAGE":      M3U8AssocLanguage,
	"FORCED":              M3U8Forced,
	"INSTREAM-ID":         M3U8InStreamId,
	"CHARACTERISTICS":     M3U8Characteristics,
	"HDCP-LEVEL":          M3U8HdcpLevel,
	"VIDEO":               M3U8Video,
	"SUBTITLES":           M3U8Subtitles,
	"CLOSED-CAPTIONS":     M3U8ClosedCaptions,
	"TIME-OFFSET":         M3U8TimeOffset,
	"PRECISE":             M3U8Precise,
	"DATA-ID":             M3U8DataId,
	"VALUE":               M3U8Value,
	"#":                   INTUnknownAttr,
	"programDataTime":     INTProgramDateTime,
	"mediaSequenceNumber": INTMediaSequenceNumber,
	"partNumber":          INTPartNumber,
}

type TagMeta struct {
	tag       TagId
	openTypes []ValueType
	attrs     []AttrId
}

var tagMeta = [...]TagMeta{
	{tag: M3U8FormatIdentifier, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXVersion, openTypes: []ValueType{valueDecimalInt}, attrs: nil},
	{tag: M3U8ExtXIndependentSegments, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXMedia, openTypes: []ValueType{}, attrs: []AttrId{M3U8Type, M3U8Uri, M3U8GroupId, M3U8Language,
		M3U8AssocLanguage, M3U8Name, M3U8Default, M3U8AutoSelect, M3U8Forced, M3U8InStreamId,
		M3U8Characteristics, M3U8Channels}},
	{tag: M3U8ExtXStreamInf, openTypes: nil, attrs: []AttrId{M3U8Bandwidth, M3U8AverageBandwidth, M3U8Codecs,
		M3U8Resolution, M3U8FrameRate, M3U8HdcpLevel, M3U8Audio, M3U8Video, M3U8Subtitles, M3U8ClosedCaptions}},
	{tag: M3U8TargetDuration, openTypes: []ValueType{valueDecimalInt}, attrs: nil},
	{tag: M3U8ExtXServerControl, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXPartInf, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXMediaSequence, openTypes: []ValueType{valueDecimalInt}, attrs: nil},
	{tag: M3U8XSkip, openTypes: nil, attrs: nil},
	{tag: M3U8ExtInf, openTypes: []ValueType{valueDecimalInt | valueUnSignedDecimalFloat, valueOptional | valueUTF8Text}, attrs: nil},
	{tag: M3U8ExtXIProgramDateTime, openTypes: []ValueType{valueDateTime}, attrs: nil},
	{tag: M3U8ExtXPart, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXPreLoadHint, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXRenditionReport, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXMap, openTypes: nil, attrs: []AttrId{M3U8Uri, M3U8ByteRange}},
	{tag: M3U8ExtXIFrameStreamInf, openTypes: nil, attrs: []AttrId{M3U8Uri}},
	{tag: M3U8ExtXDiscontinuity, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXEndList, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXPlaylistType, openTypes: []ValueType{valueEnumeratedString}, attrs: nil},
	{tag: M3U8ExtXByteRange, openTypes: []ValueType{valueDecimalInt, valueOptional | valueDecimalInt}, attrs: nil},
	{tag: M3U8ExtXKey, openTypes: nil, attrs: []AttrId{M3U8Method, M3U8IV, M3U8KeyFormat, M3U8KeyFormatVersions}},
	{tag: M3U8ExtXDataRange, openTypes: nil, attrs: []AttrId{M3U8Id}},
	{tag: M3U8ExtXDiscontinuitySequence, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXIFramesOnly, openTypes: nil, attrs: nil},
	{tag: M3U8ExtXSesionKey, openTypes: nil, attrs: []AttrId{M3U8Method, M3U8IV, M3U8KeyFormat, M3U8KeyFormatVersions}},
	{tag: M3U8ExtXSessionData, openTypes: nil, attrs: []AttrId{M3U8DataId, M3U8Value, M3U8Uri, M3U8Language}},
	{tag: M3U8ExtXStart, openTypes: nil, attrs: []AttrId{M3U8TimeOffset, M3U8Precise}},
}

type AttrMeta struct {
	attr  AttrId
	types []ValueType
}

var attrMeta = [...]AttrMeta{
	{attr: M3U8Bandwidth, types: []ValueType{valueDecimalInt}},
	{attr: M3U8AverageBandwidth, types: []ValueType{valueDecimalInt}},
	{attr: M3U8Resolution, types: []ValueType{valueDecimalResolution}},
	{attr: M3U8FrameRate, types: []ValueType{valueSignedDecimalFloat}},
	{attr: M3U8Codecs, types: []ValueType{valueQuotedString}},
	{attr: M3U8Audio, types: []ValueType{valueQuotedString}},
	{attr: M3U8Type, types: nil},
	{attr: M3U8GroupId, types: nil},
	{attr: M3U8Name, types: nil},
	{attr: M3U8Default, types: nil},
	{attr: M3U8AutoSelect, types: nil},
	{attr: M3U8Language, types: nil},
	{attr: M3U8Channels, types: nil},
	{attr: M3U8Uri, types: []ValueType{valueQuotedString}},
	{attr: M3U8CanBlockReload, types: nil},
	{attr: M3U8CanSkipUntil, types: nil},
	{attr: M3U8PartHoldBack, types: nil},
	{attr: M3U8PartTarget, types: nil},
	{attr: M3U8SkippedSegments, types: nil},
	{attr: M3U8Duration, types: nil},
	{attr: M3U8Independent, types: nil},
	{attr: M3U8LastMsn, types: nil},
	{attr: M3U8LastPart, types: nil},
	{attr: M3U8Method, types: nil},
	{attr: M3U8IV, types: nil},
	{attr: M3U8KeyFormat, types: nil},
	{attr: M3U8KeyFormatVersions, types: nil},
	{attr: M3U8ByteRange, types: nil},
	{attr: M3U8Id, types: nil},
	{attr: M3U8Class, types: nil},
	{attr: M3U8StartDate, types: nil},
	{attr: M3U8EndDate, types: nil},
	{attr: M3U8PlannedDuration, types: nil},
	{attr: M3U8Scte35Cmd, types: nil},
	{attr: M3U8Scte35Out, types: nil},
	{attr: M3U8Scte35In, types: nil},
	{attr: M3U8EndOnNext, types: nil},
	{attr: M3U8AssocLanguage, types: nil},
	{attr: M3U8Forced, types: nil},
	{attr: M3U8InStreamId, types: nil},
	{attr: M3U8Characteristics, types: nil},
	{attr: M3U8HdcpLevel, types: []ValueType{valueEnumeratedString}},
	{attr: M3U8Video, types: []ValueType{valueQuotedString}},
	{attr: M3U8Subtitles, types: []ValueType{valueQuotedString}},
	{attr: M3U8ClosedCaptions, types: []ValueType{valueQuotedString | valueEnumeratedString}},
	{attr: M3U8TimeOffset, types: nil},
	{attr: M3U8Precise, types: nil},
}

var boolToInt map[bool]int = map[bool]int{false: 0, true: 1}
