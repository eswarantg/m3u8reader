package m3u8reader

import (
	"fmt"
	"strconv"
	"time"
)

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
)

func decorateM3U8ExtXVersion(entry *M3U8Entry) error {
	if val, ok := entry.Values[m3u8UnknownKey]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8TargetDuration, val, err.Error())
		}
		entry.Values[m3u8UnknownKey] = newVal
	} else {
		return fmt.Errorf("%v missing value", M3U8TargetDuration)
	}
	return nil
}

func decorateM3U8TargetDuration(entry *M3U8Entry) error {
	if val, ok := entry.Values[m3u8UnknownKey]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8TargetDuration, val, err.Error())
		}
		entry.Values[m3u8UnknownKey] = newVal
	} else {
		return fmt.Errorf("%v missing value", M3U8TargetDuration)
	}
	return nil
}

func decorateM3U8ExtXStreamInf(entry *M3U8Entry) error {
	if _, ok := entry.Values[m3u8UnknownKey]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "URI")
	}
	if val, ok := entry.Values["BANDWIDTH"]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 32)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXStreamInf, "BANDWIDTH", val, err.Error())
		}
		entry.Values["BANDWIDTH"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "BANDWIDTH")
	}
	return nil
}

func decorateM3U8ExtXMedia(entry *M3U8Entry) error {
	if _, ok := entry.Values["URI"]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "URI")
	}
	if _, ok := entry.Values["TYPE"]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "TYPE")
	}
	if _, ok := entry.Values["LANGUAGE"]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "LANGUAGE")
	}
	if _, ok := entry.Values["GROUP-ID"]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXStreamInf, "GROUP-ID")
	}
	return nil
}

func decorateM3U8ExtInf(entry *M3U8Entry) error {
	if _, ok := entry.Values["URI"]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtInf, "URI")
	}
	if val, ok := entry.Values[m3u8UnknownKey]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8ExtInf, val, err.Error())
		}
		entry.Values[m3u8UnknownKey] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtInf, m3u8UnknownKey)
	}
	return nil
}

func decorateM3U8ExtXIProgramDateTime(entry *M3U8Entry) error {
	if val, ok := entry.Values[m3u8UnknownKey]; ok {
		newVal, err := time.Parse(time.RFC3339, val.(string))
		if err != nil {
			return fmt.Errorf("%v invalid value %v - %v", M3U8ExtXIProgramDateTime, val, err.Error())
		}
		entry.Values[m3u8UnknownKey] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXIProgramDateTime, "URI")
	}
	return nil
}

func decorateM3U8ExtXPart(entry *M3U8Entry) error {
	if _, ok := entry.Values["URI"]; !ok {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "URI")
	}
	if val, ok := entry.Values["DURATION"]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 32)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXPart, "DURATION", val, err.Error())
		}
		entry.Values["DURATION"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "DURATION")
	}
	return nil
}

func decorateM3U8ExtXMediaSequence(entry *M3U8Entry) error {
	if val, ok := entry.Values[m3u8UnknownKey]; ok {
		newVal, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXMediaSequence, m3u8UnknownKey, val, err.Error())
		}
		entry.Values[m3u8UnknownKey] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXMediaSequence, m3u8UnknownKey)
	}
	if val, ok := entry.Values["PART-TARGET"]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXMediaSequence, "PART-TARGET", val, err.Error())
		}
		entry.Values["PART-TARGET"] = newVal
	} // else failure not required it is optional
	return nil
}

func decorateM3U8ExtXPartInf(entry *M3U8Entry) error {
	if val, ok := entry.Values["PART-TARGET"]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXPartInf, "PART-TARGET", val, err.Error())
		}
		entry.Values["PART-TARGET"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPartInf, "PART-TARGET")
	}
	return nil
}

func decorateM3U8ExtXRenditionReport(entry *M3U8Entry) error {
	if val, ok := entry.Values["LAST-MSN"]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtInf, "LAST-MSN", val, err.Error())
		}
		entry.Values["LAST-MSN"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "LAST-MSN")
	}
	if val, ok := entry.Values["LAST-PART"]; ok {
		newVal, err := strconv.ParseUint(val.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtInf, "LAST-PART", val, err.Error())
		}
		entry.Values["LAST-PART"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXPart, "LAST-PART")
	}
	return nil
}

func decorateM3U8ExtXServerControl(entry *M3U8Entry) error {
	if val, ok := entry.Values["CAN-SKIP-UNTIL"]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXServerControl, "CAN-SKIP-UNTIL", val, err.Error())
		}
		entry.Values["CAN-SKIP-UNTIL"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXServerControl, "CAN-SKIP-UNTIL")
	}
	if val, ok := entry.Values["PART-HOLD-BACK"]; ok {
		newVal, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return fmt.Errorf("%v invalid value for %v = %v - %v", M3U8ExtXServerControl, "PART-HOLD-BACK", val, err.Error())
		}
		entry.Values["PART-HOLD-BACK"] = newVal
	} else {
		return fmt.Errorf("%v missing %v value", M3U8ExtXServerControl, "PART-HOLD-BACK")
	}
	return nil
}

var decorators = map[string]func(*M3U8Entry) error{
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
}
