package m3u8reader

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type TestHandler struct {
}

func (TestHandler) postRecord(tag string, kvpairs map[string]interface{}) error {
	fmt.Printf("\n%v %v", tag, kvpairs)
	return nil
}

type EmptyHandler struct {
}

func (EmptyHandler) postRecord(tag string, kvpairs map[string]interface{}) error {
	//fmt.Printf("\n%v %v", tag, kvpairs)
	return nil
}

func Test_MasterM3u8(t *testing.T) {
	//f, err := os.Open("test/ll_hls_byte_range.m3u8")
	//f, err := os.Open("test/ll_hls_delta_update.m3u8")
	//f, err := os.Open("test/ll_hls_pl.m3u8")
	//f, err := os.Open("test/index_new.m3u8")
	//f, err := os.Open("test/index_new_Variant_450k.m3u8")
	//f, err := os.Open("test/manifest.m3u8")
	f, err := os.Open("test/tv5_TS-50002_1_video.m3u8")

	if err != nil {
		t.Errorf("Unable to open file")
		return
	}
	defer f.Close()
	hdlr := TestHandler{}
	rdr := bufio.NewReader(f)
	_, err = parseM3U8_fast(rdr, hdlr)
	//_, err = parseM3U8(rdr, hdlr)
	if err != nil {
		t.Errorf("Error : %v", err)
		return
	}
}

func BenchmarkParse1(b *testing.B) {
	f, err := os.Open("test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		_, err := parseM3U8(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}

func BenchmarkParse2(b *testing.B) {
	f, err := os.Open("test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		_, err := parseM3U8_fast(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}
