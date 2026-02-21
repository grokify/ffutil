# ffutil

[![Go Reference](https://pkg.go.dev/badge/github.com/grokify/ffutil.svg)](https://pkg.go.dev/github.com/grokify/ffutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/ffutil)](https://goreportcard.com/report/github.com/grokify/ffutil)

Go wrapper for FFmpeg and FFprobe operations.

## Features

- **Type-safe command building** - Fluent API for constructing FFmpeg commands
- **Media probing** - Extract duration, resolution, codec info, and stream details
- **Hardware encoder detection** - Auto-detect VideoToolbox, NVENC, QuickSync, and other hardware encoders
- **Cross-platform** - Works on macOS, Linux, and Windows

## Installation

```bash
go get github.com/grokify/ffutil
```

Requires FFmpeg and FFprobe to be installed and available in PATH.

## Quick Start

### Probe Media Files

```go
// Get detailed media information
info, err := ffutil.Probe("video.mp4")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Duration: %v\n", info.Duration)
fmt.Printf("Resolution: %dx%d\n", info.Width, info.Height)
fmt.Printf("Video: %s, Audio: %s\n", info.VideoCodec, info.AudioCodec)

// Get just the duration
duration, err := ffutil.Duration("video.mp4")

// Check for audio/video streams
hasAudio, _ := ffutil.HasAudio("video.mp4")
hasVideo, _ := ffutil.HasVideo("video.mp4")
```

### Build FFmpeg Commands

```go
// Basic transcode
err := ffutil.New().
    Input("input.mp4").
    VideoCodec("libx264").
    AudioCodec("aac").
    Output("output.mp4").
    Run(ctx)

// With quality settings
err := ffutil.New().
    Input("input.mp4").
    VideoCodec("libx264").
    CRF(23).
    Preset("medium").
    PixelFormat("yuv420p").
    AudioCodec("aac").
    AudioBitrate("128k").
    Output("output.mp4").
    Run(ctx)

// Resize and change frame rate
err := ffutil.New().
    Input("input.mp4").
    VideoCodec("libx264").
    Size(1280, 720).
    FPS(30).
    Output("output.mp4").
    Run(ctx)

// Extract audio only
err := ffutil.New().
    Input("video.mp4").
    NoVideo().
    AudioCodec("aac").
    Output("audio.m4a").
    Run(ctx)

// Copy streams without re-encoding
err := ffutil.New().
    Input("input.mp4").
    CopyVideo().
    CopyAudio().
    Output("output.mp4").
    Run(ctx)

// Image sequence to video
err := ffutil.New().
    InputImage("frame_%04d.png", 30).
    VideoCodec("libx264").
    Duration(10.0).
    Output("output.mp4").
    Run(ctx)

// With filters
err := ffutil.New().
    Input("input.mp4").
    VideoFilter("scale=1280:720,fps=30").
    AudioFilter("volume=2.0").
    Output("output.mp4").
    Run(ctx)

// Get command string for debugging
cmd := ffutil.New().
    Input("input.mp4").
    VideoCodec("libx264").
    Output("output.mp4")
fmt.Println(cmd.String())
// Output: ffmpeg -y -i input.mp4 -c:v libx264 output.mp4
```

### Hardware Encoder Detection

```go
// Get the best available H.264 encoder
encoder := ffutil.BestH264Encoder()
fmt.Printf("Using: %s (%s)\n", encoder.Name, encoder.Type)
// macOS: h264_videotoolbox (hardware)
// Linux with NVIDIA: h264_nvenc (hardware)
// Fallback: libx264 (software)

// Check if hardware encoding is available
if ffutil.HardwareEncoderAvailable() {
    fmt.Println("Hardware acceleration available!")
}

// Use in a command
err := ffutil.New().
    Input("input.mp4").
    VideoCodec(ffutil.BestH264Encoder().Name).
    Output("output.mp4").
    Run(ctx)

// List all available encoders
encoders, _ := ffutil.ListEncoders()
for _, enc := range encoders {
    fmt.Printf("%s: %s (%s)\n", enc.Name, enc.Description, enc.Type)
}
```

### Check FFmpeg Availability

```go
// Check if both ffmpeg and ffprobe are available
if err := ffutil.Available(); err != nil {
    log.Fatalf("FFmpeg not available: %v", err)
}

// Get version strings
ffmpegVersion, _ := ffutil.Version()
ffprobeVersion, _ := ffutil.ProbeVersion()
fmt.Println(ffmpegVersion)
// Output: ffmpeg version 6.0 Copyright (c) 2000-2023...
```

## API Reference

### Command Builder Methods

| Method | Description |
|--------|-------------|
| `Input(path)` | Add input file |
| `InputWithFormat(path, format)` | Add input with format hint |
| `InputImage(path, fps)` | Add image input with loop |
| `Output(path)` | Set output file |
| `VideoCodec(codec)` | Set video codec (e.g., "libx264") |
| `AudioCodec(codec)` | Set audio codec (e.g., "aac") |
| `CopyVideo()` | Copy video stream |
| `CopyAudio()` | Copy audio stream |
| `NoVideo()` | Remove video stream |
| `NoAudio()` | Remove audio stream |
| `Size(w, h)` | Set output resolution |
| `FPS(fps)` | Set frame rate |
| `CRF(crf)` | Set quality (0-51) |
| `Preset(preset)` | Set encoding preset |
| `PixelFormat(fmt)` | Set pixel format |
| `VideoBitrate(rate)` | Set video bitrate |
| `AudioBitrate(rate)` | Set audio bitrate |
| `AudioRate(hz)` | Set audio sample rate |
| `Channels(n)` | Set audio channels |
| `Duration(sec)` | Limit output duration |
| `VideoFilter(filter)` | Set video filter |
| `AudioFilter(filter)` | Set audio filter |
| `FilterComplex(filter)` | Set complex filter |
| `Metadata(key, val)` | Set metadata |
| `Args(args...)` | Add extra arguments |
| `Build()` | Get command arguments |
| `String()` | Get full command string |
| `Run(ctx)` | Execute command |

### Probe Functions

| Function | Description |
|----------|-------------|
| `Probe(path)` | Get full media info |
| `Duration(path)` | Get duration only |
| `Resolution(path)` | Get video dimensions |
| `HasAudio(path)` | Check for audio stream |
| `HasVideo(path)` | Check for video stream |

### Encoder Functions

| Function | Description |
|----------|-------------|
| `BestH264Encoder()` | Get best available H.264 encoder |
| `BestHEVCEncoder()` | Get best available HEVC encoder |
| `EncoderAvailable(name)` | Check if encoder exists |
| `ListEncoders()` | List all video encoders |
| `HardwareEncoderAvailable()` | Check for hardware acceleration |

## License

MIT License
