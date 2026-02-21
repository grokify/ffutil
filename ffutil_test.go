package ffutil

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if !FFmpegAvailable() {
		t.Skip("ffmpeg not available")
	}

	version, err := Version()
	if err != nil {
		t.Fatalf("Version() error: %v", err)
	}

	if version == "" {
		t.Error("Version() returned empty string")
	}

	// Version should contain "ffmpeg"
	if len(version) < 6 {
		t.Errorf("Version() too short: %s", version)
	}
}

func TestProbeVersion(t *testing.T) {
	if !FFprobeAvailable() {
		t.Skip("ffprobe not available")
	}

	version, err := ProbeVersion()
	if err != nil {
		t.Fatalf("ProbeVersion() error: %v", err)
	}

	if version == "" {
		t.Error("ProbeVersion() returned empty string")
	}
}

func TestAvailable(t *testing.T) {
	err := Available()
	// We just verify it doesn't panic; the result depends on system config
	_ = err
}

func TestFFmpegAvailable(t *testing.T) {
	// Just verify it returns a boolean and doesn't panic
	result := FFmpegAvailable()
	_ = result
}

func TestFFprobeAvailable(t *testing.T) {
	// Just verify it returns a boolean and doesn't panic
	result := FFprobeAvailable()
	_ = result
}
