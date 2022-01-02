package m3u8reader_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/eswarantg/m3u8reader"
)

func Test_M3u81(t *testing.T) {
	tests := []string{
		"test/ll_hls_byte_range.m3u8",
		"test/ll_hls_delta_update.m3u8",
		"test/ll_hls_pl.m3u8",
		"test/index_new.m3u8",
		"test/index_new_Variant_450k.m3u8",
		"test/tv5.m3u8",
		"test/tv5_TS-50002_1_video.m3u8",
		"test/master.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserScanner1)
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Print(manifest.String())
	}
}

func Test_M3u82(t *testing.T) {
	tests := []string{
		"test/ll_hls_byte_range.m3u8",
		"test/ll_hls_delta_update.m3u8",
		"test/ll_hls_pl.m3u8",
		"test/index_new.m3u8",
		"test/index_new_Variant_450k.m3u8",
		"test/tv5.m3u8",
		"test/tv5_TS-50002_1_video.m3u8",
		"test/master.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserScanner2)
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Print(manifest.String())
	}
}

func Test_M3u83(t *testing.T) {
	tests := []string{
		//"test/ll_hls_byte_range.m3u8",
		//"test/ll_hls_delta_update.m3u8",
		//"test/ll_hls_pl.m3u8",
		//"test/index_new.m3u8",
		//"test/index_new_Variant_450k.m3u8",
		//"test/tv5.m3u8",
		//"test/tv5_TS-50002_1_video.m3u8",
		"test/master.m3u8",
		"test/sub.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserYacc)
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Print(manifest.String())
	}
}

func Test_M3u84(t *testing.T) {
	tests := []string{
		//"test/ll_hls_byte_range.m3u8",
		//"test/ll_hls_delta_update.m3u8",
		//"test/ll_hls_pl.m3u8",
		//"test/index_new.m3u8",
		//"test/index_new_Variant_450k.m3u8",
		//"test/tv5.m3u8",
		//"test/tv5_TS-50002_1_video.m3u8",
		"test/master.m3u8",
		"test/sub.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserGrammar)
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Print(manifest.String())
	}
}

func Test_PreloadHintEntry(t *testing.T) {
	tests := []string{
		"test/submani.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Printf("\n preloadHintEntry is: \n")
		fmt.Print(manifest.PreloadHintEntry())
	}
}

func Test_MPL(t *testing.T) {
	tests := []string{
		//"test/tv5.m3u8",
		//"test/master.m3u8",
		"test/main-manifest.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		entry, err := manifest.GetVideoMediaPlaylist(2519767)
		if err != nil {
			fmt.Printf("\nErr %v", err)
		}
		if entry != nil {
			fmt.Printf("\n%v", entry.String())
			url, err := entry.URI()
			fmt.Printf("\n%v %v", url, err)
		}
	}
}

func Test_ProgramTime(t *testing.T) {
	tests := []string{
		"test/index_new_Variant_450k.m3u8",
		"test/submani.m3u8",
		"test/sub_hls.m3u8",
	}
	for i, file := range tests {
		fmt.Printf("\n********* Test %v - %v ************", i, file)
		f, err := os.Open(file)
		if err != nil {
			t.Errorf("Unable to open file")
			return
		}
		defer f.Close()
		manifest := m3u8reader.M3U8{}
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		fmt.Print(manifest.String())
	}
}

func Benchmark_Read1(b *testing.B) {
	f, err := os.Open("test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	rdr := bytes.NewReader(manifest)
	for n := 0; n < b.N; n++ {
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserScanner1)
		_, err = manifest.Read(rdr)
		if err != nil {
			b.Errorf(err.Error())
			return
		}
	}
}

func Benchmark_Read2(b *testing.B) {
	f, err := os.Open("test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	rdr := bytes.NewReader(manifest)
	for n := 0; n < b.N; n++ {
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserScanner2)
		_, err = manifest.Read(rdr)
		if err != nil {
			b.Errorf(err.Error())
			return
		}
	}
}

func Benchmark_Read3(b *testing.B) {
	f, err := os.Open("test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	rdr := bytes.NewReader(manifest)
	for n := 0; n < b.N; n++ {
		manifest := m3u8reader.M3U8{}
		manifest.SetParserOption(m3u8reader.M3U8ParserYacc)
		_, err = manifest.Read(rdr)
		if err != nil {
			b.Errorf(err.Error())
			return
		}
	}
}
