package util

import (
	"testing"
)

func BenchmarkSimpleJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SimpleJoin("dir", "name")
	}
}

func TestGetDirectoryDepth(t *testing.T) {
	tests := []struct {
		root     string
		current  string
		expected int
	}{
		{"/home/user/project", "/home/user/project", 0},                        // Same directory
		{"/home/user/project", "/home/user/project/src", 1},                    // 1 level deep
		{"/home/user/project", "/home/user/project/src/utils", 2},              // 2 levels deep
		{"/home/user/project", "/home/user/project/src/utils/helpers", 3},      // 3 levels deep
		{"/home/user/project", "/home/user/project/src/utils/helpers/john", 4}, // 4 levels deep
		{"/", "/home", 1}, // Absolute path from root
		{"/", "/home/user", 2},
		{"/", "/home/user/project", 3},
	}

	for _, tt := range tests {
		t.Run(tt.current, func(t *testing.T) {
			got := GetDirectoryDepth(tt.root, tt.current)
			if got != tt.expected {
				t.Errorf("getDirectoryDepth(%q, %q) = %d; want %d", tt.root, tt.current, got, tt.expected)
			}
		})
	}
}

func BenchmarkGetDirectoryDepth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetDirectoryDepth("/home/user/project", "/home/user/project/src/utils/helpers/john")
	}
}

func TestFindCommonSubstrings(t *testing.T) {
	paths := []string{
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs\\node_modules",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other\\nested",
		"c:\\users\\john\\src\\github.com\\foobar\\project\\.next",
		"c:\\programs\\steam\\steamapps\\common\\Game1",
		"c:\\programs\\steam\\steamapps\\common\\Game2",
		"c:\\programs\\steam\\steamapps\\common\\Game3",
		"c:\\programs\\steam\\steamapps\\common\\Game4",
		"c:\\programs\\steam\\steamapps\\common\\Game5",
		"c:\\foo\\bar\\app1",
		"c:\\foo\\bar\\app2",
	}

	result := mapCommonPath(paths)

	expectedResult := map[string]string{
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium":                                                "c:\\users\\john\\src\\github.com",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer":                       "c:\\users\\john\\src\\github.com",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs":               "c:\\users\\john\\src\\github.com",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs\\node_modules": "c:\\users\\john\\src\\github.com",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other":                                                   "c:\\users\\john\\src\\github.com",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other\\nested":                                           "c:\\users\\john\\src\\github.com",
		"c:\\users\\john\\src\\github.com\\foobar\\project\\.next":                                                                  "c:\\users\\john\\src\\github.com",
		"c:\\programs\\steam\\steamapps\\common\\Game1":                                                                             "c:\\programs\\steam\\steamapps\\common",
		"c:\\programs\\steam\\steamapps\\common\\Game2":                                                                             "c:\\programs\\steam\\steamapps\\common",
		"c:\\programs\\steam\\steamapps\\common\\Game3":                                                                             "c:\\programs\\steam\\steamapps\\common",
		"c:\\programs\\steam\\steamapps\\common\\Game4":                                                                             "c:\\programs\\steam\\steamapps\\common",
		"c:\\programs\\steam\\steamapps\\common\\Game5":                                                                             "c:\\programs\\steam\\steamapps\\common",
		"c:\\foo\\bar\\app1": "c:\\foo\\bar",
		"c:\\foo\\bar\\app2": "c:\\foo\\bar",
	}

	for k, v := range expectedResult {
		if result[k] != v {
			t.Errorf("Expected %q, got %q", v, result[k])
		}
	}
}

func TestFindCommonSubstrings2(t *testing.T) {
	paths := []string{
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs\\node_modules",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other\\nested",
		"\\users\\john\\src\\github.com\\foobar\\project\\.next",
		"\\programs\\steam\\steamapps\\common\\Game1",
		"\\programs\\steam\\steamapps\\common\\Game2",
		"\\programs\\steam\\steamapps\\common\\Game3",
		"\\programs\\steam\\steamapps\\common\\Game4",
		"\\programs\\steam\\steamapps\\common\\Game5",
		"\\foo\\bar\\app1",
		"\\foo\\bar\\app2",
	}

	result := mapCommonPath(paths)

	expectedResult := map[string]string{
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium":                                                "\\users\\john\\src\\github.com",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer":                       "\\users\\john\\src\\github.com",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs":               "\\users\\john\\src\\github.com",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs\\node_modules": "\\users\\john\\src\\github.com",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other":                                                   "\\users\\john\\src\\github.com",
		"\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other\\nested":                                           "\\users\\john\\src\\github.com",
		"\\users\\john\\src\\github.com\\foobar\\project\\.next":                                                                  "\\users\\john\\src\\github.com",
		"\\programs\\steam\\steamapps\\common\\Game1":                                                                             "\\programs\\steam\\steamapps\\common",
		"\\programs\\steam\\steamapps\\common\\Game2":                                                                             "\\programs\\steam\\steamapps\\common",
		"\\programs\\steam\\steamapps\\common\\Game3":                                                                             "\\programs\\steam\\steamapps\\common",
		"\\programs\\steam\\steamapps\\common\\Game4":                                                                             "\\programs\\steam\\steamapps\\common",
		"\\programs\\steam\\steamapps\\common\\Game5":                                                                             "\\programs\\steam\\steamapps\\common",
		"\\foo\\bar\\app1": "\\foo\\bar",
		"\\foo\\bar\\app2": "\\foo\\bar",
	}

	for k, v := range expectedResult {
		if result[k] != v {
			t.Errorf("Expected %q, got %q", v, result[k])
		}
	}
}

func BenchmarkMapCommonPath(b *testing.B) {
	paths := []string{
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\chromium\\chromium-v112.0.0-layer\\nodejs\\node_modules",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other",
		"c:\\users\\john\\src\\github.com\\johndoe\\go\\service\\iac\\src\\other\\nested",
		"c:\\users\\john\\src\\github.com\\foobar\\project\\.next",
		"c:\\programs\\steam\\steamapps\\common\\Game1",
		"c:\\programs\\steam\\steamapps\\common\\Game2",
		"c:\\programs\\steam\\steamapps\\common\\Game3",
		"c:\\programs\\steam\\steamapps\\common\\Game4",
		"c:\\programs\\steam\\steamapps\\common\\Game5",
	}
	for i := 0; i < b.N; i++ {
		_ = mapCommonPath(paths)
	}
}
