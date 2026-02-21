# Release Notes: v0.1.0

**Release Date:** 2026-02-21

## Overview

ffutil v0.1.0 is the initial release of a Go wrapper for FFmpeg and FFprobe operations. It provides type-safe command building, media probing, and automatic hardware encoder detection across platforms.

## Highlights

- **Type-Safe Command Builder** - Fluent API for constructing FFmpeg commands
- **Media Probing** - Extract duration, resolution, codec info, and stream details
- **Hardware Encoder Detection** - Auto-detect VideoToolbox, NVENC, QuickSync, and other accelerators

## Features

### Command Builder

Build FFmpeg commands with a fluent, type-safe API:

```go
err := ffutil.New().
    Input("input.mp4").
    VideoCodec("libx264").
    CRF(23).
    Preset("medium").
    AudioCodec("aac").
    AudioBitrate("128k").
    Output("output.mp4").
    Run(ctx)
```

### Media Probing

Extract detailed information from media files:

```go
info, err := ffutil.Probe("video.mp4")
fmt.Printf("Duration: %v\n", info.Duration)
fmt.Printf("Resolution: %dx%d\n", info.Width, info.Height)
fmt.Printf("Video: %s, Audio: %s\n", info.VideoCodec, info.AudioCodec)
```

Convenience functions for common queries:

| Function | Description |
|----------|-------------|
| `Probe(path)` | Get full media info |
| `Duration(path)` | Get duration only |
| `Resolution(path)` | Get video dimensions |
| `HasAudio(path)` | Check for audio stream |
| `HasVideo(path)` | Check for video stream |

### Hardware Encoder Detection

Automatically select the best available encoder for the platform:

```go
encoder := ffutil.BestH264Encoder()
// macOS: h264_videotoolbox (hardware)
// Linux/NVIDIA: h264_nvenc (hardware)
// Fallback: libx264 (software)

err := ffutil.New().
    Input("input.mp4").
    VideoCodec(encoder.Name).
    Output("output.mp4").
    Run(ctx)
```

Supported hardware encoders:

| Platform | H.264 | HEVC |
|----------|-------|------|
| macOS | h264_videotoolbox | hevc_videotoolbox |
| NVIDIA | h264_nvenc | hevc_nvenc |
| Intel | h264_qsv | hevc_qsv |
| AMD | h264_amf | hevc_amf |
| Linux VA-API | h264_vaapi | hevc_vaapi |

## Installation

```bash
go get github.com/grokify/ffutil
```

### Prerequisites

- **FFmpeg** and **FFprobe** must be installed and available in PATH

## API Summary

### Command Builder Methods

| Method | Description |
|--------|-------------|
| `Input(path)` | Add input file |
| `Output(path)` | Set output file |
| `VideoCodec(codec)` | Set video codec |
| `AudioCodec(codec)` | Set audio codec |
| `CopyVideo()` / `CopyAudio()` | Copy streams without re-encoding |
| `NoVideo()` / `NoAudio()` | Remove streams |
| `Size(w, h)` | Set output resolution |
| `FPS(fps)` | Set frame rate |
| `CRF(crf)` | Set quality (0-51) |
| `Preset(preset)` | Set encoding preset |
| `VideoFilter(filter)` | Set video filter |
| `Run(ctx)` | Execute command |

### Encoder Functions

| Function | Description |
|----------|-------------|
| `BestH264Encoder()` | Get best available H.264 encoder |
| `BestHEVCEncoder()` | Get best available HEVC encoder |
| `EncoderAvailable(name)` | Check if encoder exists |
| `ListEncoders()` | List all video encoders |
| `HardwareEncoderAvailable()` | Check for hardware acceleration |

### Availability Functions

| Function | Description |
|----------|-------------|
| `Version()` | Get FFmpeg version string |
| `ProbeVersion()` | Get FFprobe version string |
| `Available()` | Check if both are installed |
| `FFmpegAvailable()` | Check FFmpeg only |
| `FFprobeAvailable()` | Check FFprobe only |

## Links

- [GitHub Repository](https://github.com/grokify/ffutil)
- [Go Package Documentation](https://pkg.go.dev/github.com/grokify/ffutil)
- [Changelog](https://github.com/grokify/ffutil/blob/main/CHANGELOG.md)
