package ffutil

import (
	"testing"
)

func TestProbeNonExistentFile(t *testing.T) {
	_, err := Probe("/nonexistent/file.mp4")
	if err == nil {
		t.Error("Probe() should return error for non-existent file")
	}
}

func TestDurationNonExistentFile(t *testing.T) {
	_, err := Duration("/nonexistent/file.mp4")
	if err == nil {
		t.Error("Duration() should return error for non-existent file")
	}
}

func TestResolutionNonExistentFile(t *testing.T) {
	_, _, err := Resolution("/nonexistent/file.mp4")
	if err == nil {
		t.Error("Resolution() should return error for non-existent file")
	}
}

func TestHasAudioNonExistentFile(t *testing.T) {
	_, err := HasAudio("/nonexistent/file.mp4")
	if err == nil {
		t.Error("HasAudio() should return error for non-existent file")
	}
}

func TestHasVideoNonExistentFile(t *testing.T) {
	_, err := HasVideo("/nonexistent/file.mp4")
	if err == nil {
		t.Error("HasVideo() should return error for non-existent file")
	}
}

func TestMediaInfoStruct(t *testing.T) {
	// Test that MediaInfo struct can be created
	info := &MediaInfo{
		Path:       "/test/video.mp4",
		Format:     "mp4",
		Width:      1920,
		Height:     1080,
		VideoCodec: "h264",
		AudioCodec: "aac",
		SampleRate: 44100,
		Channels:   2,
		Bitrate:    5000000,
		HasVideo:   true,
		HasAudio:   true,
	}

	if info.Path != "/test/video.mp4" {
		t.Errorf("MediaInfo.Path = %s, want /test/video.mp4", info.Path)
	}

	if info.Width != 1920 {
		t.Errorf("MediaInfo.Width = %d, want 1920", info.Width)
	}

	if info.Height != 1080 {
		t.Errorf("MediaInfo.Height = %d, want 1080", info.Height)
	}

	if !info.HasVideo {
		t.Error("MediaInfo.HasVideo should be true")
	}

	if !info.HasAudio {
		t.Error("MediaInfo.HasAudio should be true")
	}
}
