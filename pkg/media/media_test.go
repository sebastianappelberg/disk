package media

import (
	"testing"
)

func TestIsMediaFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"video.mp4", true},
		{"movie.mkv", true},
		{"clip.avi", true},
		{"film.mov", true},
		{"stream.flv", true},
		{"document.txt", false}, // Not a media file
		{"audio.mp3", true},
		{"song.wav", true},
		{"track.flac", true},
		{"image.jpg", false},        // Not a media file
		{"presentation.ppt", false}, // Not a media file
		{"archive.zip", false},      // Not a media file
		{"no_extension", false},     // No file extension
		{"upper.WMV", true},         // Case-sensitive test
		{"strange..mp4", true},      // Double dot in filename
		{"hiddenfile.", false},      // Hidden file with no extension
		{".config", false},          // Dotfile with no recognized extension
		{"video.Mp4", false},        // Mixed-case not supported.
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := isMediaFile(tt.filename)
			if got != tt.want {
				t.Errorf("isMediaFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func BenchmarkIsMediaFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = isMediaFile("media.mp4")
	}
}
