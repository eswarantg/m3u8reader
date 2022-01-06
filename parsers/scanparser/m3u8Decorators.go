package scanparser

import (
	"fmt"
	"strconv"
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
		common.M3U8Uri, common.M3U8Type,
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
	err = convertToFloat64(kv, attrs, tagId, true)
	return
}

func decorateM3U8ExtXRenditionReport(kv parsers.AttrKVPairs) (err error) {
	tagId := common.M3U8ExtXRenditionReport
	attrs := []common.AttrId{common.M3U8LastMsn, common.M3U8LastPart}
	err = convertToInt64(kv, attrs, tagId, false)
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
}

func decorateEntry(tag common.TagId, kv parsers.AttrKVPairs) (err error) {
	if decorateFn, ok := decorators[tag]; ok {
		err = decorateFn(kv)
	}
	return
}

func checkExists(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId) error {
	for _, attr := range attrIds {
		if val := kv.Get(attr); val == nil {
			return fmt.Errorf("%v missing %v value", tagId, common.AttrNames[attr])
		}
	}
	return nil
}
func convertToFloat64(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			v, ok := val.(string)
			if !ok {
				panic(fmt.Sprintf("\nconvertToFloat64 %v:%v is %T(\"%v\") not string", common.TagNames[tagId], common.AttrNames[attrId], val, val))
			}
			newVal, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("%v invalid Float value %v=\"%v\" - %v", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing Float value %v[\"%v\"]", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
func convertToInt64(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			v, ok := val.(string)
			if !ok {
				panic(fmt.Sprintf("\nconvertToInt64 %v:%v is %T(\"%v\") not string", common.TagNames[tagId], common.AttrNames[attrId], val, val))
			}
			newVal, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%v invalid Intvalue %v=\"%v\" - %v", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing Int value %v[%v", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}
func convertToTime(kv parsers.AttrKVPairs, attrIds []common.AttrId, tagId common.TagId, optional bool) error {
	for _, attrId := range attrIds {
		if val := kv.Get(attrId); val != nil {
			v, ok := val.(string)
			if !ok {
				panic(fmt.Sprintf("\nconvertToTime %v:%v is %T(\"%v\") not string", common.TagNames[tagId], common.AttrNames[attrId], val, val))
			}
			newVal, err := time.Parse(time.RFC3339Nano, v)
			if err != nil {
				return fmt.Errorf("%v invalid Time value %v=\"%v\" - %v", common.TagNames[tagId],
					common.AttrNames[attrId], val, err.Error())
			}
			kv.Store(attrId, newVal)
		} else if !optional {
			return fmt.Errorf("missing Time value %v[%v", common.TagNames[tagId], common.AttrNames[attrId])
		}
	}
	return nil
}