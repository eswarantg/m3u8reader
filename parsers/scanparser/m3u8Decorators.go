package scanparser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
)

func decorateM3U8ExtXVersion(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.INTUnknownAttr]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", common.M3U8ExtXVersion, val, err.Error())
		}
		kv[common.INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing value", common.M3U8ExtXVersion)
	}
	return nil
}

func decorateM3U8TargetDuration(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.INTUnknownAttr]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", common.M3U8TargetDuration, val, err.Error())
		}
		kv[common.INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing value", common.M3U8TargetDuration)
	}
	return nil
}

func decorateM3U8ExtXStreamInf(kv parsers.AttrKVPairs) error {
	if _, ok := kv[common.INTUnknownAttr]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXStreamInf, "URI")
	}
	if val, ok := kv[common.M3U8Bandwidth]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXStreamInf, "BANDWIDTH", val, err.Error())
		}
		kv[common.M3U8Bandwidth] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXStreamInf, "BANDWIDTH")
	}
	return nil
}

func decorateM3U8ExtXMedia(kv parsers.AttrKVPairs) error {
	if _, ok := kv[common.M3U8Uri]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXStreamInf, "URI")
	}
	if _, ok := kv[common.M3U8Type]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXStreamInf, "TYPE")
	}
	if _, ok := kv[common.M3U8Language]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXStreamInf, "LANGUAGE")
	}
	if _, ok := kv[common.M3U8GroupId]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXStreamInf, "GROUP-ID")
	}
	return nil
}

func decorateM3U8ExtInf(kv parsers.AttrKVPairs) error {
	if _, ok := kv[common.M3U8Uri]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtInf, "URI")
	}
	if val, ok := kv[common.INTUnknownAttr]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", common.M3U8ExtInf, val, err.Error())
		}
		kv[common.INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtInf, common.INTUnknownAttr)
	}
	return nil
}

func decorateM3U8ExtXIProgramDateTime(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.INTUnknownAttr]; ok {
		newVal, err := time.Parse(time.RFC3339Nano, val.(string))
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", common.M3U8ExtXIProgramDateTime, val, err.Error())
		}
		kv[common.INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXIProgramDateTime, "URI")
	}
	return nil
}

func decorateM3U8ExtXPart(kv parsers.AttrKVPairs) error {
	if _, ok := kv[common.M3U8Uri]; !ok {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXPart, "URI")
	}
	if val, ok := kv[common.M3U8Duration]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXPart, "DURATION", val, err.Error())
		}
		kv[common.M3U8Duration] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXPart, "DURATION")
	}
	return nil
}

func decorateM3U8ExtXMediaSequence(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.INTUnknownAttr]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXMediaSequence, common.INTUnknownAttr, val, err.Error())
		}
		kv[common.INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXMediaSequence, common.INTUnknownAttr)
	}
	if val, ok := kv[common.M3U8PartTarget]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXMediaSequence, "PART-TARGET", val, err.Error())
		}
		kv[common.M3U8PartTarget] = newVal
	} // else failure not required it is optional
	return nil
}

func decorateM3U8ExtXPartInf(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.M3U8PartTarget]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXPartInf, "PART-TARGET", val, err.Error())
		}
		kv[common.M3U8PartTarget] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXPartInf, "PART-TARGET")
	}
	return nil
}

func decorateM3U8ExtXRenditionReport(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.M3U8LastMsn]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtInf, "LAST-MSN", val, err.Error())
		}
		kv[common.M3U8LastMsn] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXPart, "LAST-MSN")
	}
	if val, ok := kv[common.M3U8LastPart]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtInf, "LAST-PART", val, err.Error())
		}
		kv[common.M3U8LastPart] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXPart, "LAST-PART")
	}
	return nil
}

func decorateM3U8ExtXServerControl(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.M3U8CanSkipUntil]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXServerControl, "CAN-SKIP-UNTIL", val, err.Error())
		}
		kv[common.M3U8CanSkipUntil] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXServerControl, "CAN-SKIP-UNTIL")
	}
	if val, ok := kv[common.M3U8PartHoldBack]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8ExtXServerControl, "PART-HOLD-BACK", val, err.Error())
		}
		kv[common.M3U8PartHoldBack] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8ExtXServerControl, "PART-HOLD-BACK")
	}
	return nil
}

func decorateM3U8XSkip(kv parsers.AttrKVPairs) error {
	if val, ok := kv[common.M3U8SkippedSegments]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", common.M3U8XSkip, "SKIPPED-SEGMENTS", val, err.Error())
		}
		kv[common.M3U8SkippedSegments] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", common.M3U8XSkip, "SKIPPED-SEGMENTS")
	}
	return nil
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
