package m3u8reader_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/eswarantg/m3u8reader"
)

func Test_M3u8(t *testing.T) {
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
		fmt.Print(manifest.String())
		fmt.Printf("\n preloadHintEntry is: \n")
		fmt.Print(manifest.PreloadHintEntry())
	}
}

func Test_MPL(t *testing.T) {
	tests := []string{
		//"test/tv5.m3u8",
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
		_, err = manifest.Read(f)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		entry, err := manifest.GetVideoMediaPlaylist(10000000)
		if entry != nil {
			fmt.Printf("\n%v %v", entry.String(), err)
		}
	}
}
