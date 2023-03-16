package common

//All TagNames in M3U8
var TagNames = [...]string{
	"UNKNOWNTAG",
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
	M3U8UNKNOWNTAG TagId = iota
	M3U8FormatIdentifier
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

var TagToTagId map[string]TagId = map[string]TagId{
	"UNKNOWNTAG":                   M3U8UNKNOWNTAG,
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

var AttrNames = [...]string{
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
	"TITLE",
	"BYTERANGE-START",
	"BYTERANGE-LENGTH",
	"#",
	"programDataTime",
	"mediaSequenceNumber",
	"partNumber",
	"PROGRAM-ID",
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
	M3U8Title
	M3U8ByteRangeStart
	M3U8ByteRangeLength
	INTUnknownAttr
	INTProgramDateTime
	INTMediaSequenceNumber
	INTPartNumber
	M3U8ProgramId
)

var AttrToAttrId map[string]AttrId = map[string]AttrId{
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
	"TITLE":               M3U8Title,
	"BYTERANGE-START":     M3U8ByteRangeStart,
	"BYTERANGE-LENGTH":    M3U8ByteRangeLength,
	"#":                   INTUnknownAttr,
	"programDataTime":     INTProgramDateTime,
	"mediaSequenceNumber": INTMediaSequenceNumber,
	"partNumber":          INTPartNumber,
	"PROGRAM-ID":          M3U8ProgramId,
}
