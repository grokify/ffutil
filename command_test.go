package ffutil

import (
	"strings"
	"testing"
)

func TestCommandBuild(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *Command
		contains    []string
		notContains []string
	}{
		{
			name: "basic transcode",
			cmd: New().
				Input("input.mp4").
				VideoCodec("libx264").
				AudioCodec("aac").
				Output("output.mp4"),
			contains: []string{
				"-y", "-i", "input.mp4",
				"-c:v", "libx264",
				"-c:a", "aac",
				"output.mp4",
			},
		},
		{
			name: "copy streams",
			cmd: New().
				Input("input.mp4").
				CopyVideo().
				CopyAudio().
				Output("output.mp4"),
			contains: []string{
				"-c:v", "copy",
				"-c:a", "copy",
			},
		},
		{
			name: "no audio",
			cmd: New().
				Input("input.mp4").
				VideoCodec("libx264").
				NoAudio().
				Output("output.mp4"),
			contains:    []string{"-an"},
			notContains: []string{"-c:a"},
		},
		{
			name: "no video",
			cmd: New().
				Input("input.mp4").
				AudioCodec("aac").
				NoVideo().
				Output("output.mp4"),
			contains:    []string{"-vn"},
			notContains: []string{"-c:v"},
		},
		{
			name: "with size and fps",
			cmd: New().
				Input("input.mp4").
				VideoCodec("libx264").
				Size(1920, 1080).
				FPS(30).
				Output("output.mp4"),
			contains: []string{
				"-s", "1920x1080",
				"-r", "30",
			},
		},
		{
			name: "with crf and preset",
			cmd: New().
				Input("input.mp4").
				VideoCodec("libx264").
				CRF(23).
				Preset("medium").
				Output("output.mp4"),
			contains: []string{
				"-crf", "23",
				"-preset", "medium",
			},
		},
		{
			name: "with pixel format",
			cmd: New().
				Input("input.mp4").
				VideoCodec("libx264").
				PixelFormat("yuv420p").
				Output("output.mp4"),
			contains: []string{"-pix_fmt", "yuv420p"},
		},
		{
			name: "with bitrates",
			cmd: New().
				Input("input.mp4").
				VideoCodec("libx264").
				VideoBitrate("5M").
				AudioCodec("aac").
				AudioBitrate("128k").
				Output("output.mp4"),
			contains: []string{
				"-b:v", "5M",
				"-b:a", "128k",
			},
		},
		{
			name: "with audio settings",
			cmd: New().
				Input("input.mp4").
				AudioCodec("aac").
				AudioRate(44100).
				Channels(2).
				Output("output.mp4"),
			contains: []string{
				"-ar", "44100",
				"-ac", "2",
			},
		},
		{
			name: "with video filter",
			cmd: New().
				Input("input.mp4").
				VideoFilter("scale=1280:720").
				Output("output.mp4"),
			contains: []string{"-vf", "scale=1280:720"},
		},
		{
			name: "with audio filter",
			cmd: New().
				Input("input.mp4").
				AudioFilter("volume=2.0").
				Output("output.mp4"),
			contains: []string{"-af", "volume=2.0"},
		},
		{
			name: "with complex filter",
			cmd: New().
				Input("input.mp4").
				FilterComplex("[0:v]scale=1280:720[v]").
				Output("output.mp4"),
			contains: []string{"-filter_complex", "[0:v]scale=1280:720[v]"},
		},
		{
			name: "with duration",
			cmd: New().
				Input("input.mp4").
				Duration(10.5).
				Output("output.mp4"),
			contains: []string{"-t", "10.500"},
		},
		{
			name: "with metadata",
			cmd: New().
				Input("input.mp4").
				Metadata("title", "My Video").
				Output("output.mp4"),
			contains: []string{"-metadata", "title=My Video"},
		},
		{
			name: "image input with loop",
			cmd: New().
				InputImage("image.png", 30).
				Duration(5.0).
				VideoCodec("libx264").
				Output("output.mp4"),
			contains: []string{
				"-loop", "1",
				"-framerate", "30",
				"-i", "image.png",
			},
		},
		{
			name: "input with format",
			cmd: New().
				InputWithFormat("input.raw", "rawvideo").
				Output("output.mp4"),
			contains: []string{"-f", "rawvideo"},
		},
		{
			name: "extra args",
			cmd: New().
				Input("input.mp4").
				Args("-movflags", "+faststart").
				Output("output.mp4"),
			contains: []string{"-movflags", "+faststart"},
		},
		{
			name: "no overwrite",
			cmd: New().
				Input("input.mp4").
				Overwrite(false).
				Output("output.mp4"),
			notContains: []string{"-y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.cmd.Build()
			argsStr := strings.Join(args, " ")

			for _, want := range tt.contains {
				if !strings.Contains(argsStr, want) {
					t.Errorf("Build() missing %q in %v", want, args)
				}
			}

			for _, notWant := range tt.notContains {
				if strings.Contains(argsStr, notWant) {
					t.Errorf("Build() should not contain %q in %v", notWant, args)
				}
			}
		})
	}
}

func TestCommandString(t *testing.T) {
	cmd := New().
		Input("input.mp4").
		VideoCodec("libx264").
		Output("output.mp4")

	str := cmd.String()

	if !strings.HasPrefix(str, "ffmpeg ") {
		t.Errorf("String() should start with 'ffmpeg ', got: %s", str)
	}

	if !strings.Contains(str, "input.mp4") {
		t.Errorf("String() should contain input file, got: %s", str)
	}
}

func TestCommandStringWithSpaces(t *testing.T) {
	cmd := New().
		Input("my video.mp4").
		Output("my output.mp4")

	str := cmd.String()

	// Files with spaces should be quoted
	if !strings.Contains(str, `"my video.mp4"`) {
		t.Errorf("String() should quote paths with spaces, got: %s", str)
	}
}

func TestMultipleInputs(t *testing.T) {
	cmd := New().
		Input("video.mp4").
		Input("audio.mp3").
		Output("output.mp4")

	args := cmd.Build()
	argsStr := strings.Join(args, " ")

	// Should have two -i flags
	count := strings.Count(argsStr, "-i ")
	if count != 2 {
		t.Errorf("Build() should have 2 inputs, got %d", count)
	}
}
