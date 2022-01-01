package m3u8reader

import (
	"fmt"
	"strconv"
	"time"
)

func decorateM3U8ExtXVersion(entry *M3U8Entry) error {
	if val, ok := entry.Values[INTUnknownAttr]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8ExtXVersion, val, err.Error())
		}
		entry.Values[INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing value", M3U8ExtXVersion)
	}
	return nil
}

func decorateM3U8TargetDuration(entry *M3U8Entry) error {
	if val, ok := entry.Values[INTUnknownAttr]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8TargetDuration, val, err.Error())
		}
		entry.Values[INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing value", M3U8TargetDuration)
	}
	return nil
}

func decorateM3U8ExtXStreamInf(entry *M3U8Entry) error {
	if _, ok := entry.Values[INTUnknownAttr]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "URI")
	}
	if val, ok := entry.Values[M3U8Bandwidth]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXStreamInf, "BANDWIDTH", val, err.Error())
		}
		entry.Values[M3U8Bandwidth] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "BANDWIDTH")
	}
	return nil
}

func decorateM3U8ExtXMedia(entry *M3U8Entry) error {
	if _, ok := entry.Values[M3U8Uri]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "URI")
	}
	if _, ok := entry.Values[M3U8Type]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "TYPE")
	}
	if _, ok := entry.Values[M3U8Language]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "LANGUAGE")
	}
	if _, ok := entry.Values[M3U8GroupId]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "GROUP-ID")
	}
	return nil
}

func decorateM3U8ExtInf(entry *M3U8Entry) error {
	if _, ok := entry.Values[M3U8Uri]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtInf, "URI")
	}
	if val, ok := entry.Values[INTUnknownAttr]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8ExtInf, val, err.Error())
		}
		entry.Values[INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtInf, m3u8UnknownKey)
	}
	return nil
}

func decorateM3U8ExtXIProgramDateTime(entry *M3U8Entry) error {
	if val, ok := entry.Values[INTUnknownAttr]; ok {
		newVal, err := time.Parse(time.RFC3339, val.(string))
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8ExtXIProgramDateTime, val, err.Error())
		}
		entry.Values[INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXIProgramDateTime, "URI")
	}
	return nil
}

func decorateM3U8ExtXPart(entry *M3U8Entry) error {
	if _, ok := entry.Values[M3U8Uri]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "URI")
	}
	if val, ok := entry.Values[M3U8Duration]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXPart, "DURATION", val, err.Error())
		}
		entry.Values[M3U8Duration] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "DURATION")
	}
	return nil
}

func decorateM3U8ExtXMediaSequence(entry *M3U8Entry) error {
	if val, ok := entry.Values[INTUnknownAttr]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXMediaSequence, m3u8UnknownKey, val, err.Error())
		}
		entry.Values[INTUnknownAttr] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXMediaSequence, m3u8UnknownKey)
	}
	if val, ok := entry.Values[M3U8PartTarget]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXMediaSequence, "PART-TARGET", val, err.Error())
		}
		entry.Values[M3U8PartTarget] = newVal
	} // else failure not required it is optional
	return nil
}

func decorateM3U8ExtXPartInf(entry *M3U8Entry) error {
	if val, ok := entry.Values[M3U8PartTarget]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXPartInf, "PART-TARGET", val, err.Error())
		}
		entry.Values[M3U8PartTarget] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPartInf, "PART-TARGET")
	}
	return nil
}

func decorateM3U8ExtXRenditionReport(entry *M3U8Entry) error {
	if val, ok := entry.Values[M3U8LastMsn]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtInf, "LAST-MSN", val, err.Error())
		}
		entry.Values[M3U8LastMsn] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "LAST-MSN")
	}
	if val, ok := entry.Values[M3U8LastPart]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtInf, "LAST-PART", val, err.Error())
		}
		entry.Values[M3U8LastPart] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "LAST-PART")
	}
	return nil
}

func decorateM3U8ExtXServerControl(entry *M3U8Entry) error {
	if val, ok := entry.Values[M3U8CanSkipUntil]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXServerControl, "CAN-SKIP-UNTIL", val, err.Error())
		}
		entry.Values[M3U8CanSkipUntil] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXServerControl, "CAN-SKIP-UNTIL")
	}
	if val, ok := entry.Values[M3U8PartHoldBack]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXServerControl, "PART-HOLD-BACK", val, err.Error())
		}
		entry.Values[M3U8PartHoldBack] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXServerControl, "PART-HOLD-BACK")
	}
	return nil
}

func decorateM3U8XSkip(entry *M3U8Entry) error {
	if val, ok := entry.Values[M3U8SkippedSegments]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8XSkip, "SKIPPED-SEGMENTS", val, err.Error())
		}
		entry.Values[M3U8SkippedSegments] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8XSkip, "SKIPPED-SEGMENTS")
	}
	return nil
}

var decorators = map[TagId]func(*M3U8Entry) error{
	M3U8ExtXVersion:          decorateM3U8ExtXVersion,
	M3U8TargetDuration:       decorateM3U8TargetDuration,
	M3U8ExtXStreamInf:        decorateM3U8ExtXStreamInf,
	M3U8ExtXMedia:            decorateM3U8ExtXMedia,
	M3U8ExtInf:               decorateM3U8ExtInf,
	M3U8ExtXIProgramDateTime: decorateM3U8ExtXIProgramDateTime,
	M3U8ExtXPart:             decorateM3U8ExtXPart,
	M3U8ExtXMediaSequence:    decorateM3U8ExtXMediaSequence,
	M3U8ExtXPartInf:          decorateM3U8ExtXPartInf,
	M3U8ExtXRenditionReport:  decorateM3U8ExtXRenditionReport,
	M3U8ExtXServerControl:    decorateM3U8ExtXServerControl,
	M3U8XSkip:                decorateM3U8XSkip,
}
