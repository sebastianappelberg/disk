// This file contains code from the project github.com/razsteinmetz/go-ptn,
// with modifications to fit the rest of the project.
// Source: https://github.com/razsteinmetz/go-ptn.
// License: MIT
package torrents

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

// test routine using the testdata.json file
const testdata = "testdata.json"

type TestData struct {
	Name     string      `json:"name"`
	Expected TorrentInfo `json:"expected"`
}

func TestParser(t *testing.T) {
	var testData []TestData
	buf, err := os.ReadFile(testdata)
	if err != nil {
		t.Fatalf("error reading golden filke: %v", err)
	}
	err = json.Unmarshal(buf, &testData)
	if err != nil {
		t.Fatalf("error trying to unmarshal the test data: %v", err)
	}
	for i, data := range testData {
		name := data.Name
		expected := data.Expected
		t.Run(fmt.Sprintf("Testing: %s", name), func(t *testing.T) {
			actual, err := ParseName(name)
			if err != nil {
				t.Fatalf("test %v - %s: parser error:\n  %v", i, name, err)
			}
			if !reflect.DeepEqual(*actual, expected) {
				t.Fatalf("test %v: wrong result for %q\nwant:\n  %v\ngot:\n  %v", i, name, expected, *actual)
			}
		})
	}
}

func TestParserSpecific(t *testing.T) {
	data := TestData{
		Name: "Archer (2009) Season 13 S13 (1080p AMZN WEB-DL x265 HEVC 10bit EAC3 5.1 Ghost)",
		Expected: TorrentInfo{
			Title:   "Life's Too Short",
			Season:  1,
			Episode: 0,
			IsMovie: false,
		},
	}
	name := data.Name
	expected := data.Expected
	t.Run(fmt.Sprintf("Testing: %s", name), func(t *testing.T) {
		actual, err := ParseName(name)
		if err != nil {
			t.Fatalf("test - %s: parser error:\n  %v", name, err)
		}
		if !reflect.DeepEqual(*actual, expected) {
			t.Fatalf("test: wrong result for %q\nwant:\n  %v\ngot:\n  %v", name, expected, *actual)
		}
	})
}
