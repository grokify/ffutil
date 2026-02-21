package ffutil

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Command represents an ffmpeg command builder.
type Command struct {
	inputs       []inputSpec
	outputPath   string
	videoCodec   string
	audioCodec   string
	videoBitrate string
	audioBitrate string
	width        int
	height       int
	fps          int
	crf          int
	preset       string
	pixelFormat  string
	audioRate    int
	channels     int
	duration     float64
	startTime    float64
	copyVideo    bool
	copyAudio    bool
	noAudio      bool
	noVideo      bool
	overwrite    bool
	extraArgs    []string
	filterVideo  string
	filterAudio  string
	filterComplex string
	metadata     map[string]string
}

// inputSpec represents an input file with optional parameters.
type inputSpec struct {
	path      string
	format    string
	frameRate int
	loop      bool
	duration  float64
	startTime float64
}

// New creates a new FFmpeg command builder.
func New() *Command {
	return &Command{
		overwrite: true,
		metadata:  make(map[string]string),
	}
}

// Input adds an input file to the command.
func (c *Command) Input(path string) *Command {
	c.inputs = append(c.inputs, inputSpec{path: path})
	return c
}

// InputWithFormat adds an input file with a specified format.
func (c *Command) InputWithFormat(path, format string) *Command {
	c.inputs = append(c.inputs, inputSpec{path: path, format: format})
	return c
}

// InputImage adds an image input with loop enabled.
func (c *Command) InputImage(path string, frameRate int) *Command {
	c.inputs = append(c.inputs, inputSpec{
		path:      path,
		loop:      true,
		frameRate: frameRate,
	})
	return c
}

// InputWithDuration adds an input with a duration limit.
func (c *Command) InputWithDuration(path string, duration float64) *Command {
	c.inputs = append(c.inputs, inputSpec{path: path, duration: duration})
	return c
}

// Output sets the output file path.
func (c *Command) Output(path string) *Command {
	c.outputPath = path
	return c
}

// VideoCodec sets the video codec (e.g., "libx264", "h264_videotoolbox").
func (c *Command) VideoCodec(codec string) *Command {
	c.videoCodec = codec
	c.copyVideo = false
	return c
}

// AudioCodec sets the audio codec (e.g., "aac", "libmp3lame").
func (c *Command) AudioCodec(codec string) *Command {
	c.audioCodec = codec
	c.copyAudio = false
	return c
}

// CopyVideo copies the video stream without re-encoding.
func (c *Command) CopyVideo() *Command {
	c.copyVideo = true
	c.videoCodec = ""
	return c
}

// CopyAudio copies the audio stream without re-encoding.
func (c *Command) CopyAudio() *Command {
	c.copyAudio = true
	c.audioCodec = ""
	return c
}

// NoAudio removes the audio stream from output.
func (c *Command) NoAudio() *Command {
	c.noAudio = true
	return c
}

// NoVideo removes the video stream from output.
func (c *Command) NoVideo() *Command {
	c.noVideo = true
	return c
}

// Size sets the output video dimensions.
func (c *Command) Size(width, height int) *Command {
	c.width = width
	c.height = height
	return c
}

// FPS sets the output frame rate.
func (c *Command) FPS(fps int) *Command {
	c.fps = fps
	return c
}

// CRF sets the Constant Rate Factor for quality (0-51, lower is better).
func (c *Command) CRF(crf int) *Command {
	c.crf = crf
	return c
}

// Preset sets the encoding preset (e.g., "ultrafast", "medium", "slow").
func (c *Command) Preset(preset string) *Command {
	c.preset = preset
	return c
}

// PixelFormat sets the pixel format (e.g., "yuv420p").
func (c *Command) PixelFormat(format string) *Command {
	c.pixelFormat = format
	return c
}

// VideoBitrate sets the video bitrate (e.g., "5M", "2000k").
func (c *Command) VideoBitrate(bitrate string) *Command {
	c.videoBitrate = bitrate
	return c
}

// AudioBitrate sets the audio bitrate (e.g., "128k", "320k").
func (c *Command) AudioBitrate(bitrate string) *Command {
	c.audioBitrate = bitrate
	return c
}

// AudioRate sets the audio sample rate in Hz.
func (c *Command) AudioRate(rate int) *Command {
	c.audioRate = rate
	return c
}

// Channels sets the number of audio channels.
func (c *Command) Channels(channels int) *Command {
	c.channels = channels
	return c
}

// Duration limits the output duration in seconds.
func (c *Command) Duration(seconds float64) *Command {
	c.duration = seconds
	return c
}

// StartTime sets the start time for processing.
func (c *Command) StartTime(seconds float64) *Command {
	c.startTime = seconds
	return c
}

// Overwrite enables or disables overwriting output files.
func (c *Command) Overwrite(overwrite bool) *Command {
	c.overwrite = overwrite
	return c
}

// VideoFilter sets the video filter graph.
func (c *Command) VideoFilter(filter string) *Command {
	c.filterVideo = filter
	return c
}

// AudioFilter sets the audio filter graph.
func (c *Command) AudioFilter(filter string) *Command {
	c.filterAudio = filter
	return c
}

// FilterComplex sets a complex filter graph.
func (c *Command) FilterComplex(filter string) *Command {
	c.filterComplex = filter
	return c
}

// Metadata sets a metadata key-value pair.
func (c *Command) Metadata(key, value string) *Command {
	c.metadata[key] = value
	return c
}

// Args adds extra arguments to the command.
func (c *Command) Args(args ...string) *Command {
	c.extraArgs = append(c.extraArgs, args...)
	return c
}

// Build returns the ffmpeg command arguments.
func (c *Command) Build() []string {
	var args []string

	// Global options
	if c.overwrite {
		args = append(args, "-y")
	}

	// Input options
	for _, input := range c.inputs {
		if input.loop {
			args = append(args, "-loop", "1")
		}
		if input.frameRate > 0 {
			args = append(args, "-framerate", strconv.Itoa(input.frameRate))
		}
		if input.format != "" {
			args = append(args, "-f", input.format)
		}
		if input.duration > 0 {
			args = append(args, "-t", formatDuration(input.duration))
		}
		if input.startTime > 0 {
			args = append(args, "-ss", formatDuration(input.startTime))
		}
		args = append(args, "-i", input.path)
	}

	// Filter options
	if c.filterComplex != "" {
		args = append(args, "-filter_complex", c.filterComplex)
	}
	if c.filterVideo != "" {
		args = append(args, "-vf", c.filterVideo)
	}
	if c.filterAudio != "" {
		args = append(args, "-af", c.filterAudio)
	}

	// Video options
	if c.noVideo {
		args = append(args, "-vn")
	} else if c.copyVideo {
		args = append(args, "-c:v", "copy")
	} else if c.videoCodec != "" {
		args = append(args, "-c:v", c.videoCodec)
	}

	if c.width > 0 && c.height > 0 {
		args = append(args, "-s", fmt.Sprintf("%dx%d", c.width, c.height))
	}

	if c.fps > 0 {
		args = append(args, "-r", strconv.Itoa(c.fps))
	}

	if c.crf > 0 {
		args = append(args, "-crf", strconv.Itoa(c.crf))
	}

	if c.preset != "" {
		args = append(args, "-preset", c.preset)
	}

	if c.pixelFormat != "" {
		args = append(args, "-pix_fmt", c.pixelFormat)
	}

	if c.videoBitrate != "" {
		args = append(args, "-b:v", c.videoBitrate)
	}

	// Audio options
	if c.noAudio {
		args = append(args, "-an")
	} else if c.copyAudio {
		args = append(args, "-c:a", "copy")
	} else if c.audioCodec != "" {
		args = append(args, "-c:a", c.audioCodec)
	}

	if c.audioBitrate != "" {
		args = append(args, "-b:a", c.audioBitrate)
	}

	if c.audioRate > 0 {
		args = append(args, "-ar", strconv.Itoa(c.audioRate))
	}

	if c.channels > 0 {
		args = append(args, "-ac", strconv.Itoa(c.channels))
	}

	// Output options
	if c.duration > 0 {
		args = append(args, "-t", formatDuration(c.duration))
	}

	if c.startTime > 0 {
		args = append(args, "-ss", formatDuration(c.startTime))
	}

	// Metadata
	for key, value := range c.metadata {
		args = append(args, "-metadata", fmt.Sprintf("%s=%s", key, value))
	}

	// Extra arguments
	args = append(args, c.extraArgs...)

	// Output path
	if c.outputPath != "" {
		args = append(args, c.outputPath)
	}

	return args
}

// String returns the full ffmpeg command as a string.
func (c *Command) String() string {
	args := c.Build()
	quoted := make([]string, len(args))
	for i, arg := range args {
		if strings.ContainsAny(arg, " \t\n\"'") {
			quoted[i] = fmt.Sprintf("%q", arg)
		} else {
			quoted[i] = arg
		}
	}
	return "ffmpeg " + strings.Join(quoted, " ")
}

// Run executes the ffmpeg command.
func (c *Command) Run(ctx context.Context) error {
	args := c.Build()
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nstderr: %s", err, stderr.String())
	}
	return nil
}

// RunWithOutput executes the ffmpeg command and returns combined output.
func (c *Command) RunWithOutput(ctx context.Context) ([]byte, error) {
	args := c.Build()
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("ffmpeg failed: %w\noutput: %s", err, string(output))
	}
	return output, nil
}

// formatDuration formats a duration in seconds for ffmpeg.
func formatDuration(seconds float64) string {
	return strconv.FormatFloat(seconds, 'f', 3, 64)
}
