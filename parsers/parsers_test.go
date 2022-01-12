package parsers_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/eswarantg/m3u8reader/common"
	"github.com/eswarantg/m3u8reader/parsers"
	"github.com/eswarantg/m3u8reader/parsers/scanparser"
)

type TestHandler struct {
}

func (TestHandler) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
	fmt.Printf("\n%v(%v)", common.TagNames[tag], tag)
	if kvpairs != nil {
		for k, v := range kvpairs.Map() {
			fmt.Printf("\n\t%v(%v)=%v", common.AttrNames[k], k, v)
		}
		kvpairs.Done()
		fmt.Printf("\n")
	}
	return nil
}

type EmptyHandler struct {
}

func (EmptyHandler) PostRecord(tag common.TagId, kvpairs *parsers.AttrKVPairs) error {
	//fmt.Printf("\n%v %v", tag, kvpairs)
	if kvpairs != nil {
		kvpairs.Done()
	}
	return nil
}

//Helper
func ReadFile(path string, t *testing.T) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		t.Errorf("Unable to open file")
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func Test_MediaM3u8(t *testing.T) {
	var hdlr parsers.M3u8Handler
	files := [...]string{
		"../test/sub.m3u8",
	}
	buffer := make([]byte, 4096)
	for _, file := range files {
		data, err := ReadFile(file, t)
		if err != nil {
			continue
		}

		fmt.Printf("\n****** ScanParser3 - w/o sync ******")
		hdlr = TestHandler{}
		parsers.AttrKVPairsSyncPool = false
		scanner1 := scanparser.ScanParser3{}
		_, err = scanner1.ParseData(data, hdlr, buffer)
		if err != nil {
			t.Errorf("Error : %v", err)
			continue
		}
		/*
			fmt.Printf("\n****** ScanParser3 - w/ sync ******")
			hdlr = TestHandler{}
			parsers.AttrKVPairsSyncPool = true
			_, err = scanner1.ParseData(data, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
		/*
			fmt.Printf("\n****** ScanParser2 ******")
			hdlr = TestHandler{}
			scanner2 := scanparser.ScanParser2{}
			_, err = scanner2.ParseData(data, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
		/*
			fmt.Printf("****** YaccParser ******")
			hdlr = TestHandler{}
			scanner3 := yaccparser.YaccParser{}
			_, err = scanner3.ParseData(data, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
		/*
			fmt.Printf("\n****** GrammarParser ******")
			hdlr = TestHandler{}
			scanner4 := grammarparser.GrammarParser{}
			_, err = scanner4.ParseData(data, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
	}
}

func Test_MasterM3u8(t *testing.T) {
	var hdlr parsers.M3u8Handler
	var files = [...]string{
		"../test/manifest.m3u8",
	}
	buffer := make([]byte, 4096)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()

		fmt.Printf("\n****** ScanParser3 - w/o sync ******")
		hdlr = TestHandler{}
		parsers.AttrKVPairsSyncPool = false
		scanner1 := scanparser.ScanParser3{}
		_, err = scanner1.Parse(f, hdlr, buffer)
		if err != nil {
			t.Errorf("Error : %v", err)
			continue
		}

		fmt.Printf("\n****** ScanParser3 - w/ sync ******")
		hdlr = TestHandler{}
		parsers.AttrKVPairsSyncPool = true
		_, err = scanner1.Parse(f, hdlr, buffer)
		if err != nil {
			t.Errorf("Error : %v", err)
			continue
		}

		/*
			hdlr = TestHandler{}
			scanner2 := scanparser.ScanParser2{}
			f.Seek(0, io.SeekStart)
			_, err = scanner2.Parse(f, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
		/*
			hdlr = TestHandler{}
			scanner3 := yaccparser.YaccParser{}
			f.Seek(0, io.SeekStart)
			_, err = scanner3.Parse(f, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
		/*
			hdlr = TestHandler{}
			scanner4 := grammarparser.GrammarParser{}
			f.Seek(0, io.SeekStart)
			_, err = scanner4.Parse(f, hdlr)
			if err != nil {
				t.Errorf("Error : %v", err)
				continue
			}
		*/
	}
}

/*
func BenchmarkParse1a(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	parsers.AttrKVPairsSyncPool = false
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		scanner := scanparser.ScanParser1{}
		_, err = scanner.Parse(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}

func BenchmarkParse1b(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	parsers.AttrKVPairsSyncPool = true
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		scanner := scanparser.ScanParser1{}
		_, err = scanner.Parse(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}

func BenchmarkParse2a(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	parsers.AttrKVPairsSyncPool = false
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		scanner := scanparser.ScanParser3{}
		_, err = scanner.Parse(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}
*/

/*
func BenchmarkParse2b(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	parsers.AttrKVPairsSyncPool = true
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		scanner := scanparser.ScanParser3{}
		_, err = scanner.Parse(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}
*/

func BenchmarkParse2b_1(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	buffer := make([]byte, 4096)
	manifest, err := ioutil.ReadAll(f)
	if err != nil {
		b.Errorf("Error : %v", err)
	}
	f.Close()
	parsers.AttrKVPairsSyncPool = true
	for n := 0; n < b.N; n++ {
		hdlr := EmptyHandler{}
		rdr := bytes.NewReader(manifest)
		scanner := scanparser.ScanParser3{}
		_, err = scanner.Parse(rdr, hdlr, buffer)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}

/*
func BenchmarkParse2(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
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
		scanner := scanparser.ScanParser2{}
		_, err = scanner.Parse(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}
*/

/*
func BenchmarkParse3(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
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
		scanner := yaccparser.YaccParser{}
		_, err = scanner.Parse(rdr, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}
*/

/*
func BenchmarkParse4(b *testing.B) {
	f, err := os.Open("../test/sub.m3u8")
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
		scanner := grammarparser.GrammarParser{}
		_, err = scanner.ParseData(manifest, hdlr)
		if err != nil {
			b.Errorf("Error : %v", err)
			return
		}
	}
}
*/
