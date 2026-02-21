// Package ffutil provides a Go wrapper for FFmpeg and FFprobe operations.
//
// This package simplifies common video/audio operations by providing:
//   - Type-safe command building
//   - Consistent error handling
//   - Media probing (duration, resolution, codec info)
//   - Hardware encoder detection
//
// Basic usage:
//
//	// Probe media file
//	info, err := ffutil.Probe("video.mp4")
//	fmt.Println(info.Duration)
//
//	// Build and run ffmpeg command
//	err := ffutil.New().
//	    Input("input.mp4").
//	    VideoCodec("libx264").
//	    Output("output.mp4").
//	    Run(ctx)
package ffutil

import (
	"fmt"
	"os/exec"
	"strings"
)

// Version returns the ffmpeg version string.
func Version() (string, error) {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", nil
}

// ProbeVersion returns the ffprobe version string.
func ProbeVersion() (string, error) {
	cmd := exec.Command("ffprobe", "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffprobe not found: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", nil
}

// Available checks if ffmpeg and ffprobe are available in PATH.
func Available() error {
	if _, err := Version(); err != nil {
		return err
	}
	if _, err := ProbeVersion(); err != nil {
		return err
	}
	return nil
}

// FFmpegAvailable checks if ffmpeg is available in PATH.
func FFmpegAvailable() bool {
	cmd := exec.Command("ffmpeg", "-version")
	return cmd.Run() == nil
}

// FFprobeAvailable checks if ffprobe is available in PATH.
func FFprobeAvailable() bool {
	cmd := exec.Command("ffprobe", "-version")
	return cmd.Run() == nil
}
