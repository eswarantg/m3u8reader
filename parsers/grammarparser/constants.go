package grammarparser

import "github.com/eswarantg/m3u8reader/common"

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

type OpenType struct {
	types ValueType
	attr  common.AttrId
}
type TagMeta struct {
	tag       common.TagId
	openTypes []OpenType
	attrs     []common.AttrId
}

var tagMeta = [...]TagMeta{
	{tag: common.M3U8UNKNOWNTAG, openTypes: nil, attrs: nil},
	{tag: common.M3U8FormatIdentifier, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXVersion, openTypes: []OpenType{
		{types: valueDecimalInt, attr: common.INTUnknownAttr},
	}, attrs: nil},
	{tag: common.M3U8ExtXIndependentSegments, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXMedia, openTypes: nil, attrs: []common.AttrId{common.M3U8Type,
		common.M3U8Uri, common.M3U8GroupId, common.M3U8Language, common.M3U8AssocLanguage, common.M3U8Name,
		common.M3U8Default, common.M3U8AutoSelect, common.M3U8Forced, common.M3U8InStreamId,
		common.M3U8Characteristics, common.M3U8Channels}},
	{tag: common.M3U8ExtXStreamInf, openTypes: []OpenType{
		{types: valueNextLineEnumeratedString, attr: common.INTUnknownAttr},
	}, attrs: []common.AttrId{common.M3U8Bandwidth,
		common.M3U8AverageBandwidth, common.M3U8Codecs, common.M3U8Resolution, common.M3U8FrameRate,
		common.M3U8HdcpLevel, common.M3U8Audio, common.M3U8Video, common.M3U8Subtitles,
		common.M3U8ClosedCaptions}},
	{tag: common.M3U8TargetDuration, openTypes: []OpenType{
		{types: valueDecimalInt, attr: common.INTUnknownAttr},
	}, attrs: nil},
	{tag: common.M3U8ExtXServerControl, openTypes: nil, attrs: []common.AttrId{common.M3U8CanBlockReload,
		common.M3U8CanSkipUntil, common.M3U8PartHoldBack}},
	{tag: common.M3U8ExtXPartInf, openTypes: nil, attrs: []common.AttrId{common.M3U8PartTarget}},
	{tag: common.M3U8ExtXMediaSequence, openTypes: []OpenType{
		{types: valueDecimalInt, attr: common.INTUnknownAttr},
	}, attrs: nil},
	{tag: common.M3U8XSkip, openTypes: nil, attrs: []common.AttrId{common.M3U8SkippedSegments}},
	{tag: common.M3U8ExtInf, openTypes: []OpenType{
		{types: valueDecimalInt | valueUnSignedDecimalFloat, attr: common.INTUnknownAttr},
		{types: valueEnumeratedString | valueOptional, attr: common.M3U8Title},
		{types: valueNextLineEnumeratedString, attr: common.M3U8Uri},
	}, attrs: nil},
	{tag: common.M3U8ExtXIProgramDateTime, openTypes: []OpenType{
		{types: valueDateTime, attr: common.INTUnknownAttr},
	}, attrs: nil},
	{tag: common.M3U8ExtXPart, openTypes: nil, attrs: []common.AttrId{common.M3U8Duration,
		common.M3U8Independent, common.M3U8Uri}},
	{tag: common.M3U8ExtXPreLoadHint, openTypes: nil, attrs: []common.AttrId{
		common.M3U8Type, common.M3U8Uri}},
	{tag: common.M3U8ExtXRenditionReport, openTypes: nil, attrs: []common.AttrId{
		common.M3U8Uri, common.M3U8LastMsn, common.M3U8LastPart}},
	{tag: common.M3U8ExtXMap, openTypes: nil, attrs: []common.AttrId{common.M3U8Uri, common.M3U8ByteRange}},
	{tag: common.M3U8ExtXIFrameStreamInf, openTypes: nil, attrs: []common.AttrId{common.M3U8Uri}},
	{tag: common.M3U8ExtXDiscontinuity, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXEndList, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXPlaylistType, openTypes: []OpenType{
		{types: valueEnumeratedString, attr: common.INTUnknownAttr},
	}, attrs: nil},
	{tag: common.M3U8ExtXByteRange, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXKey, openTypes: nil, attrs: []common.AttrId{common.M3U8Method, common.M3U8IV,
		common.M3U8KeyFormat, common.M3U8KeyFormatVersions}},
	{tag: common.M3U8ExtXDataRange, openTypes: nil, attrs: []common.AttrId{common.M3U8Id}},
	{tag: common.M3U8ExtXDiscontinuitySequence, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXIFramesOnly, openTypes: nil, attrs: nil},
	{tag: common.M3U8ExtXSesionKey, openTypes: nil, attrs: []common.AttrId{common.M3U8Method, common.M3U8IV,
		common.M3U8KeyFormat, common.M3U8KeyFormatVersions}},
	{tag: common.M3U8ExtXSessionData, openTypes: nil, attrs: []common.AttrId{common.M3U8DataId, common.M3U8Value,
		common.M3U8Uri, common.M3U8Language}},
	{tag: common.M3U8ExtXStart, openTypes: nil, attrs: []common.AttrId{common.M3U8TimeOffset, common.M3U8Precise}},
}

type AttrMeta struct {
	attr  common.AttrId
	types []ValueType
}

//For now map to only 1 value
//If int/float => map to float
//if quoted/enumerate => map to BITOR() - as it needs different handling
var attrMeta = [...]AttrMeta{
	{attr: common.M3U8Bandwidth, types: []ValueType{valueDecimalInt}},
	{attr: common.M3U8AverageBandwidth, types: []ValueType{valueDecimalInt}},
	{attr: common.M3U8Resolution, types: []ValueType{valueDecimalResolution}},
	{attr: common.M3U8FrameRate, types: []ValueType{valueSignedDecimalFloat}},
	{attr: common.M3U8Codecs, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Audio, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Type, types: []ValueType{valueEnumeratedString}},
	{attr: common.M3U8GroupId, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Name, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Default, types: []ValueType{valueEnumeratedString}},
	{attr: common.M3U8AutoSelect, types: []ValueType{valueEnumeratedString}},
	{attr: common.M3U8Language, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Channels, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Uri, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8CanBlockReload, types: []ValueType{valueEnumeratedString}},
	{attr: common.M3U8CanSkipUntil, types: []ValueType{valueDecimalInt}},
	{attr: common.M3U8PartHoldBack, types: []ValueType{valueSignedDecimalFloat}},
	{attr: common.M3U8PartTarget, types: []ValueType{valueSignedDecimalFloat}},
	{attr: common.M3U8SkippedSegments, types: []ValueType{valueDecimalInt}},
	{attr: common.M3U8Duration, types: []ValueType{valueSignedDecimalFloat}},
	{attr: common.M3U8Independent, types: []ValueType{valueEnumeratedString}},
	{attr: common.M3U8LastMsn, types: []ValueType{valueDecimalInt}},
	{attr: common.M3U8LastPart, types: []ValueType{valueDecimalInt}},
	{attr: common.M3U8Method, types: nil},
	{attr: common.M3U8IV, types: nil},
	{attr: common.M3U8KeyFormat, types: nil},
	{attr: common.M3U8KeyFormatVersions, types: nil},
	{attr: common.M3U8ByteRange, types: nil},
	{attr: common.M3U8Id, types: nil},
	{attr: common.M3U8Class, types: nil},
	{attr: common.M3U8StartDate, types: nil},
	{attr: common.M3U8EndDate, types: nil},
	{attr: common.M3U8PlannedDuration, types: nil},
	{attr: common.M3U8Scte35Cmd, types: nil},
	{attr: common.M3U8Scte35Out, types: nil},
	{attr: common.M3U8Scte35In, types: nil},
	{attr: common.M3U8EndOnNext, types: nil},
	{attr: common.M3U8AssocLanguage, types: nil},
	{attr: common.M3U8Forced, types: nil},
	{attr: common.M3U8InStreamId, types: nil},
	{attr: common.M3U8Characteristics, types: nil},
	{attr: common.M3U8HdcpLevel, types: []ValueType{valueEnumeratedString}},
	{attr: common.M3U8Video, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8Subtitles, types: []ValueType{valueQuotedString}},
	{attr: common.M3U8ClosedCaptions, types: []ValueType{valueQuotedString | valueEnumeratedString}},
	{attr: common.M3U8TimeOffset, types: nil},
	{attr: common.M3U8Precise, types: nil},
	{attr: common.M3U8DataId, types: nil},
	{attr: common.M3U8Value, types: nil},
	{attr: common.M3U8Title, types: nil},
}
