package scanparser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

func decorateM3U8ExtXVersion(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXVersion
	attrs := []common.AttrId{common.INTUnknownAttr}
	err = convertToInt64(kv, attrs, tagId, false)
	return
}

func decorateM3U8TargetDuration(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8TargetDuration
	attrs := []common.AttrId{common.INTUnknownAttr}
	err = convertToInt64(kv, attrs, tagId, false)
	return
}

func decorateM3U8ExtXStreamInf(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXStreamInf
	attrs := []common.AttrId{common.INTUnknownAttr}
	err = checkExists(kv, attrs, tagId)
	if err != nil {
		return
	}
	attrs = []common.AttrId{common.M3U8Bandwidth}
	err = convertToInt64(kv, attrs, tagId, false)
	return err
}

func decorateM3U8ExtXMedia(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXMedia
	attrs := []common.AttrId{
		common.M3U8Type,
		common.M3U8Language, common.M3U8GroupId,
	}
	err = checkExists(kv, attrs, tagId)
	return
}

func decorateM3U8ExtInf(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtInf
	attrs := []common.AttrId{common.M3U8Uri}
	err = checkExists(kv, attrs, tagId)
	if err != nil {
		return
	}
	attrs = []common.AttrId{common.INTUnknownAttr}
	err = convertToFloat64(kv, attrs, tagId, false)
	return err
}

func decorateM3U8ExtXIProgramDateTime(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXIProgramDateTime
	attrs := []common.AttrId{common.INTUnknownAttr}
	err = convertToTime(kv, attrs, tagId, false)
	return
}

func decorateM3U8ExtXPart(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXPart
	attrs := []common.AttrId{common.M3U8Uri}
	err = checkExists(kv, attrs, tagId)
	if err != nil {
		return
	}
	attrs = []common.AttrId{common.M3U8ByteRange}
	err = convertToByteRange(kv, attrs, tagId, true) //optional
	if err != nil {
		return
	}
	attrs = []common.AttrId{common.M3U8Duration}
	err = convertToFloat64(kv, attrs, tagId, false)
	return
}

func decorateM3U8ExtXMediaSequence(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXMediaSequence
	attrs := []common.AttrId{common.INTUnknownAttr}
	err = convertToInt64(kv, attrs, tagId, false)
	return
}

func decorateM3U8ExtXPartInf(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXPartInf
	attrs := []common.AttrId{common.M3U8PartTarget}
	err = convertToFloat64(kv, attrs, tagId, true) //optional
	return
}

func decorateM3U8ExtXRenditionReport(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXRenditionReport
	attrs := []common.AttrId{common.M3U8LastMsn, common.M3U8LastPart}
	err = convertToInt64(kv, attrs, tagId, false)
	if err != nil {
		return
	}
	attrs = []common.AttrId{common.M3U8ByteRangeStart}
	err = convertToInt64(kv, attrs, tagId, true) //optional
	return
}

func decorateM3U8ExtXServerControl(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXServerControl
	attrs := []common.AttrId{common.M3U8CanSkipUntil, common.M3U8PartHoldBack}
	err = convertToFloat64(kv, attrs, tagId, false)
	return
}

func decorateM3U8XSkip(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8XSkip
	attrs := []common.AttrId{common.M3U8SkippedSegments}
	err = convertToInt64(kv, attrs, tagId, false)
	return
}

func decorateM3U8ExtXPreLoadHint(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXPreLoadHint
	attrs := []common.AttrId{common.M3U8Uri}
	err = checkExists(kv, attrs, tagId)
	if err != nil {
		return
	}
	attrs = []common.AttrId{common.M3U8ByteRangeStart, common.M3U8ByteRangeLength}
	err = convertToInt64(kv, attrs, tagId, true) //optional
	return
}

var decorators = map[common.TagId]func(kv parsers.AttrKVPairs) error{
	common.M3U8ExtXVersion:          decorateM3U8ExtXVersion,
	common.M3U8TargetDuration:       decorateM3U8TargetDuration,
	common.M3U8ExtXStreamInf:        decorateM3U8ExtXStreamInf,
	common.M3U8ExtXMedia:            decorateM3U8ExtXMedia,
	common.M3U8ExtInf:               decorateM3U8ExtInf,
	common.M3U8ExtXIProgramDateTime: decorateM3U8ExtXIProgramDateTime,
	common.M3U8ExtXPart:             decorateM3U8ExtXPart,
	common.M3U8ExtXMediaSequence:    decorateM3U8ExtXMediaSequence,
	common.M3U8ExtXPartInf:          decorateM3U8ExtXPartInf,
	common.M3U8ExtXRenditionReport:  decorateM3U8ExtXRenditionReport,
	common.M3U8ExtXServerControl:    decorateM3U8ExtXServerControl,
	common.M3U8XSkip:                decorateM3U8XSkip,
	common.M3U8ExtXPreLoadHint:      decorateM3U8ExtXPreLoadHint,
}

func decorateEntry(tag common.TagId, kv parsers.AttrKVPairs) (err error) {
	if decorateFn, ok := decorators[tag]; ok {
		err = decorateFn(kv)
	}
	return
}

func checkExists(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId) error {
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val == nil {
			return fmt.Errorf("missing \"%v\":\"%v\" value", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
func convertToFloat64(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	var newVal float64
	var err error
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			switch v := val.(type) {
			case []byte:
				newVal, err = strconv.ParseFloat(string(v), 64)
			case string:
				newVal, err = strconv.ParseFloat(v, 64)
			default:
				panic(fmt.Sprintf("\nconvertToFloat64 \"%v\":\"%v\" is %T(\"%v\") not string", common.TagNames[tagId],
					common.AttrNames[attrId], v, v))
			}
			if err != nil {
				return fmt.Errorf("invalid Float value \"%v\":\"%v\"=\"%v\" - \"%v\"", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing \"%v\":\"%v\" Float value", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
func convertToInt64(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	var newVal int64
	var err error
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			switch v := val.(type) {
			case []byte:
				newVal, err = strconv.ParseInt(string(v), 10, 64)
			case string:
				newVal, err = strconv.ParseInt(v, 10, 64)
			default:
				panic(fmt.Sprintf("\nconvertToInt64 \"%v\":\"%v\" is %T(\"%v\") not string", common.TagNames[tagId],
					common.AttrNames[attrId], v, v))
			}
			if err != nil {
				return fmt.Errorf("invalid Int value \"%v\":\"%v\"=\"%v\" - \"%v\"", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing \"%v\":\"%v\" Int value", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
func convertToByteRange(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	//Ref: https://datatracker.ietf.org/doc/html/draft-pantos-hls-rfc8216bis#section-4.4.4.2

	//	#EXT-X-BYTERANGE:<n>[@<o>]

	//   where n is a decimal-integer indicating the length of the sub-range
	//in bytes.  If present, o is a decimal-integer indicating the start of
	//the sub-range, as a byte offset from the beginning of the resource.
	//If o is not present, the sub-range begins at the next byte following
	//the sub-range of the previous Media Segment.

	//If o is not present, a previous Media Segment MUST appear in the
	//Playlist file and MUST be a sub-range of the same media resource, or
	//the Media Segment is undefined and the client MUST fail to parse the
	//Playlist.

	//https://datatracker.ietf.org/doc/html/draft-pantos-hls-rfc8216bis#section-4.4.4.9
	//EXT-X-PART
	//BYTERANGE

	//The value is a quoted-string specifying a byte range into the
	//resource identified by the URI attribute.  This range SHOULD
	//contain only the Media Initialization Section.  The format of the
	//byte range is described in Section 4.4.4.2.  This attribute is
	//OPTIONAL; if it is not present, the byte range is the entire
	//resource indicated by the URI.

	var newVal [2]int64
	var err error
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			switch v := val.(type) {
			case []byte:
				parts := bytes.Split(v, []byte{'@'})
				switch len(parts) {
				case 2:
					newVal[1], err = strconv.ParseInt(string(parts[1]), 10, 64)
					if err != nil {
						break
					}
					newVal[0], err = strconv.ParseInt(string(parts[0]), 10, 64)
				case 1:
					err = errors.New("byteRange with 1 part not supported, expect n@o")
					//newVal[1] = 0
					//newVal[0], err = strconv.ParseInt(string(parts[0]), 10, 64)
				default:
					err = errors.New("byteRange expected 1 part or 2 parts with @ seperator")
				}
			case string:
				parts := strings.Split(v, "@")
				switch len(parts) {
				case 2:
					newVal[1], err = strconv.ParseInt(parts[1], 10, 64)
					if err != nil {
						break
					}
					newVal[0], err = strconv.ParseInt(parts[0], 10, 64)
				case 1:
					newVal[1] = -1
					newVal[0], err = strconv.ParseInt(parts[0], 10, 64)
				default:
					err = errors.New("byteRange expected 1 part or 2 parts with @ seperator")
				}
			default:
				panic(fmt.Sprintf("\nconvertToByteRange \"%v\":\"%v\" is %T(\"%v\") not string", common.TagNames[tagId],
					common.AttrNames[attrId], v, v))
			}
			if err != nil {
				return fmt.Errorf("invalid byteRange value \"%v\":\"%v\"=\"%v\" - \"%v\"", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing \"%v\":\"%v\" byteRange value", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
func convertToTime(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	var newVal time.Time
	var err error
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			switch v := val.(type) {
			case []byte:
				newVal, err = time.Parse(time.RFC3339Nano, string(v))
			case string:
				newVal, err = time.Parse(time.RFC3339Nano, v)
			default:
				panic(fmt.Sprintf("\nconvertToTime %v:%v is %T(\"%v\") not string", common.TagNames[tagId],
					common.AttrNames[attrId], v, v))
			}
			if err != nil {
				return fmt.Errorf("invalid Time value \"%v\":\"%v\"=\"%v\" - \"%v\"", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing \"%v\":\"%v\" Time value", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
