// Package torrents.
// This file contains code from the project github.com/razsteinmetz/go-ptn,
// with modifications to fit the rest of the project.
// Source: https://github.com/razsteinmetz/go-ptn.
// License: MIT
package torrents

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TorrentInfo is the resulting structure returned by ParseName.
type TorrentInfo struct {
	Title      string `json:"title,omitempty"`
	Season     int    `json:"season,omitempty"`
	Episode    int    `json:"episode,omitempty"`
	Year       int    `json:"year,omitempty"`
	Resolution string `json:"resolution,omitempty"` //1080p etc
	Quality    string `json:"quality,omitempty"`
	Codec      string `json:"codec,omitempty"`
	Audio      string `json:"audio,omitempty"`
	Service    string `json:"service,omitempty"` // NF etc
	Group      string `json:"group,omitempty"`
	Region     string `json:"region,omitempty"`
	Extended   bool   `json:"extended,omitempty"`
	Hardcoded  bool   `json:"hardcoded,omitempty"`
	Limited    bool   `json:"limited,omitempty"`
	Proper     bool   `json:"proper,omitempty"`
	Repack     bool   `json:"repack,omitempty"` // also rerip
	Container  string `json:"container,omitempty"`
	Widescreen bool   `json:"widescreen,omitempty"`
	Website    string `json:"website,omitempty"`
	Language   string `json:"language,omitempty"`
	Sbs        string `json:"sbs,omitempty"`
	Unrated    bool   `json:"unrated,omitempty"`
	Size       string `json:"size,omitempty"`
	Threed     bool   `json:"3d,omitempty"`
	Country    string `json:"country,omitempty"`
	IsMovie    bool   `json:"isMovie"`
}

func setField(tor *TorrentInfo, field, val string) {
	// set the Field by reflecting its info
	ttor := reflect.TypeOf(tor)
	torV := reflect.ValueOf(tor)

	caser := cases.Title(language.English)
	field = caser.String(field)

	v, _ := ttor.Elem().FieldByName(field)
	switch v.Type.Kind() {
	case reflect.Bool:
		torV.Elem().FieldByName(field).SetBool(true)
	case reflect.Int:
		clean, _ := strconv.ParseInt(val, 10, 64)
		torV.Elem().FieldByName(field).SetInt(clean)
	case reflect.Uint:
		clean, _ := strconv.ParseUint(val, 10, 64)
		torV.Elem().FieldByName(field).SetUint(clean)
	case reflect.String:
		torV.Elem().FieldByName(field).SetString(val)
	default:
		// Do nothing.
	}
}

// ParseName breaks up the given filename into TorrentInfo.
func ParseName(filename string) (*TorrentInfo, error) {
	tor := &TorrentInfo{}
	var startIndex, endIndex = 0, len(filename)
	// remove any underline and replace with Spaces
	cleanName := strings.Replace(filename, "_", " ", -1)
	if matches := container.FindAllStringSubmatch(cleanName, -1); len(matches) != 0 {
		tor.Container = matches[0][1]
		cleanName = cleanName[0 : len(cleanName)-4]
	} else if matches := otherExtensions.FindAllStringSubmatch(cleanName, -1); len(matches) != 0 {
		cleanName = cleanName[0 : len(cleanName)-4] // remove the . and the extension from the checked strings.
	}
	// go over all patterns
	for _, pattern := range patterns {
		matches := pattern.re.FindAllStringSubmatch(cleanName, -1)
		if len(matches) == 0 {
			continue
		}
		matchIdx := 0
		if pattern.last {
			// Take last occurrence of element.
			matchIdx = len(matches) - 1
		}

		index := strings.Index(cleanName, matches[matchIdx][1])
		if index == 0 {
			startIndex = len(matches[matchIdx][1])
		} else if index < endIndex {
			endIndex = index
		}
		if startIndex > endIndex {
			endIndex = len(filename)
			continue
		}
		setField(tor, pattern.name, matches[matchIdx][2])
	}

	// Start process for title and remove all dots/underscore from it
	raw := strings.Split(filename[startIndex:endIndex], "(")[0]
	cleanName = raw
	if strings.HasPrefix(cleanName, "- ") {
		cleanName = raw[2:]
	}
	// clean out the title remove any starting chars
	cleanName = strings.Trim(cleanName, " -_.^/\\(){}[]")
	// only remove the dots if there are no spaces for some titles have dots and spaces
	if strings.ContainsRune(cleanName, '.') && !strings.ContainsRune(cleanName, ' ') {
		cleanName = strings.Replace(cleanName, ".", " ", -1)
	}
	cleanName = strings.ReplaceAll(cleanName, "_", " ")
	cleanName = strings.ReplaceAll(cleanName, "  ", " ")
	cleanName = strings.Trim(cleanName, " -_.^/\\(){}[]")
	if matches := countryre.FindAllStringSubmatch(cleanName, -1); len(matches) != 0 {
		tor.Country = matches[0][1]
		cleanName = cleanName[0 : len(cleanName)-3] // Remove the country from the name.
	}
	setField(tor, "title", cleanName)
	tor.IsMovie = tor.Episode == 0 && tor.Season == 0
	return tor, nil
}

// this is the file extension, we only consider extensions that are relevant to movies but srt are possibles, but they
// are not containers
var container = regexp.MustCompile(`(?i)\.(MKV|AVI|MP4|MOV|MPG|MPEG|FLV|F4V|SWF|WMV|MP2|MPE|MPV|OGG|M4V|M4P|AVCHD)$`)
var otherExtensions = regexp.MustCompile(`(?i)\.(SRT|SUB|IDX)$`) // not containers but reasonable to expect
var countryre = regexp.MustCompile(`\b(US|UK)$`)                 // country name in upper case at the end of the title

var patterns = []struct {
	name string
	// Use the last matching pattern. E.g. Year.
	last bool
	kind reflect.Kind
	// REs need to have 2 sub expressions (groups), the first one is "raw", and
	// the second one for the "clean" value.
	// E.g. Episode matching on "S01E18" will result in: raw = "E18", clean = "18".
	re *regexp.Regexp
}{
	{"season", false, reflect.Int, regexp.MustCompile(`(?i)(s?([0-9]{1,2}))[ex]`)},
	{"episode", false, reflect.Int, regexp.MustCompile(`(?i)([ex]([0-9]{2})(?:[^0-9]|$))`)},
	{"year", true, reflect.Int, regexp.MustCompile(`\b(((?:19[0-9]|20[0-9])[0-9]))\b`)},
	{"resolution", false, reflect.String, regexp.MustCompile(`\b(([0-9]{3,4}p))\b`)},
	{"quality", false, reflect.String, regexp.MustCompile(`(?i)\b(((?:PPV\.)?[HP]DTV|(?:HD)?CAM|B[DR]Rip|(?:HD-?)?TS|(?:PPV )?WEB-?DL(?: DVDRip)?|HDRip|DVDRip|DVDRIP|CamRip|W[EB]BRip|BluRay|DvDScr|telesync))\b`)},
	{"codec", false, reflect.String, regexp.MustCompile(`(?i)\b((xvid|exvid|[hx]\.?26[45]))\b`)},
	{"audio", false, reflect.String, regexp.MustCompile(`(?i)\b((MP3|DD5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC[.-]LC|AAC(?:\.?2\.0)?|AC3(?:\.5\.1)?))\b`)},
	{"region", false, reflect.String, regexp.MustCompile(`(?i)\b(R([0-9]))\b`)},
	{"size", false, reflect.String, regexp.MustCompile(`(?i)\b((\d+(?:\.\d+)?(?:GB|MB)))\b`)},
	{"website", false, reflect.String, regexp.MustCompile(`^(\[ ?([^\]]+?) ?\])`)},
	{"language", false, reflect.String, regexp.MustCompile(`(?i)\b((rus\.eng|ita\.eng))\b`)},
	{"sbs", false, reflect.String, regexp.MustCompile(`(?i)\b(((?:Half-)?SBS))\b`)},
	{"group", false, reflect.String, regexp.MustCompile(`\b(- ?([^-]+(?:-={[^-]+-?$)?))$`)},
	{"service", false, reflect.String, regexp.MustCompile(`(?i)\b((NF|ATVP|BBC|CBS|ABC|DSNP|DSNY|FOX|HMAX|HULU|iP|MTV|NICK|SHO))\b`)},
	{"extended", false, reflect.Bool, regexp.MustCompile(`(?i)\b(EXTENDED(:?.CUT)?)\b`)},
	{"hardcoded", false, reflect.Bool, regexp.MustCompile(`(?i)\b((HC))\b`)},
	{"proper", false, reflect.Bool, regexp.MustCompile(`(?i)\b((PROPER))\b`)},
	{"repack", false, reflect.Bool, regexp.MustCompile(`(?i)\b((REPACK|RERIP))\b`)},
	{"limited", false, reflect.Bool, regexp.MustCompile(`(?i)\b((LIMITED))\b`)},
	{"widescreen", false, reflect.Bool, regexp.MustCompile(`(?i)\b((WS))\b`)},
	{"unrated", false, reflect.Bool, regexp.MustCompile(`(?i)\b((UNRATED))\b`)},
	{"threeD", false, reflect.Bool, regexp.MustCompile(`(?i)\b((3D))\b`)},
}

// sanity check the patterns for capture groups
func init() {
	for _, pat := range patterns {
		if pat.re.NumSubexp() != 2 {
			fmt.Printf("Pattern %q does not have enough capture groups. want 2, got %d\n", pat.name, pat.re.NumSubexp())
			os.Exit(1)
		}
	}
}
