package ffutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// MediaInfo contains information about a media file.
type MediaInfo struct {
	// Path is the file path that was probed
	Path string `json:"path"`

	// Duration is the media duration
	Duration time.Duration `json:"duration"`

	// Format is the container format (e.g., "mp4", "webm", "mp3")
	Format string `json:"format"`

	// Width is the video width in pixels (0 if no video)
	Width int `json:"width,omitempty"`

	// Height is the video height in pixels (0 if no video)
	Height int `json:"height,omitempty"`

	// VideoCodec is the video codec name (empty if no video)
	VideoCodec string `json:"videoCodec,omitempty"`

	// AudioCodec is the audio codec name (empty if no audio)
	AudioCodec string `json:"audioCodec,omitempty"`

	// SampleRate is the audio sample rate in Hz (0 if no audio)
	SampleRate int `json:"sampleRate,omitempty"`

	// Channels is the number of audio channels (0 if no audio)
	Channels int `json:"channels,omitempty"`

	// Bitrate is the overall bitrate in bits per second
	Bitrate int64 `json:"bitrate,omitempty"`

	// HasVideo indicates if the file has a video stream
	HasVideo bool `json:"hasVideo"`

	// HasAudio indicates if the file has an audio stream
	HasAudio bool `json:"hasAudio"`
}

// ffprobeOutput represents the JSON output from ffprobe
type ffprobeOutput struct {
	Format  ffprobeFormat   `json:"format"`
	Streams []ffprobeStream `json:"streams"`
}

type ffprobeFormat struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	BitRate    string `json:"bit_rate"`
}

type ffprobeStream struct {
	CodecType  string `json:"codec_type"`
	CodecName  string `json:"codec_name"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	SampleRate string `json:"sample_rate,omitempty"`
	Channels   int    `json:"channels,omitempty"`
}

// Probe returns detailed information about a media file.
func Probe(path string) (*MediaInfo, error) {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	}

	cmd := exec.Command("ffprobe", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w (stderr: %s)", err, stderr.String())
	}

	var output ffprobeOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	info := &MediaInfo{
		Path:   path,
		Format: output.Format.FormatName,
	}

	// Parse duration
	if output.Format.Duration != "" {
		if dur, err := strconv.ParseFloat(output.Format.Duration, 64); err == nil {
			info.Duration = time.Duration(dur * float64(time.Second))
		}
	}

	// Parse bitrate
	if output.Format.BitRate != "" {
		if br, err := strconv.ParseInt(output.Format.BitRate, 10, 64); err == nil {
			info.Bitrate = br
		}
	}

	// Process streams
	for _, stream := range output.Streams {
		switch stream.CodecType {
		case "video":
			info.HasVideo = true
			info.VideoCodec = stream.CodecName
			info.Width = stream.Width
			info.Height = stream.Height
		case "audio":
			info.HasAudio = true
			info.AudioCodec = stream.CodecName
			info.Channels = stream.Channels
			if stream.SampleRate != "" {
				if sr, err := strconv.Atoi(stream.SampleRate); err == nil {
					info.SampleRate = sr
				}
			}
		}
	}

	return info, nil
}

// Duration returns the duration of a media file.
// This is a convenience function that only fetches duration.
func Duration(path string) (time.Duration, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	}

	cmd := exec.Command("ffprobe", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w (stderr: %s)", err, stderr.String())
	}

	durStr := strings.TrimSpace(stdout.String())
	if durStr == "" || durStr == "N/A" {
		return 0, fmt.Errorf("duration not available for %s", path)
	}

	dur, err := strconv.ParseFloat(durStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration %q: %w", durStr, err)
	}

	return time.Duration(dur * float64(time.Second)), nil
}

// Resolution returns the video resolution (width, height) of a media file.
// Returns (0, 0) if the file has no video stream.
func Resolution(path string) (width, height int, err error) {
	info, err := Probe(path)
	if err != nil {
		return 0, 0, err
	}
	return info.Width, info.Height, nil
}

// HasAudio returns true if the media file has an audio stream.
func HasAudio(path string) (bool, error) {
	info, err := Probe(path)
	if err != nil {
		return false, err
	}
	return info.HasAudio, nil
}

// HasVideo returns true if the media file has a video stream.
func HasVideo(path string) (bool, error) {
	info, err := Probe(path)
	if err != nil {
		return false, err
	}
	return info.HasVideo, nil
}
